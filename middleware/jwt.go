package middleware

import (
	"fama-api/database/models"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	"regexp"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// JwtPayload contains the values stored in the jwt Token used by the app
type JwtPayload struct {
	UserID   uint
	Username string
	IsAdmin  bool
}

//GetJwtMiddleware returns the Jwt middleware used by the app
func GetJwtMiddleware(db *gorm.DB) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(os.Getenv("SECRET_KEY")),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: "id",
		Authenticator: func(c *gin.Context) (interface{}, error) {
			type LoginPayload struct {
				Username string `form:"username" binding:"required"`
				Password string `form:"password" binding:"required"`
			}

			var payload LoginPayload

			if err := c.Bind(&payload); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			var user struct {
				ID       uint
				Username string
				Password string
				IsAdmin  bool
			}
			if err := db.Table("users").
				Select("users.id, users.username, users.password, CASE WHEN admins.user_id IS NULL THEN 'false' ELSE 'true' END as is_admin").
				Where(&models.User{Username: payload.Username}).
				Joins("LEFT JOIN admins ON admins.user_id = users.id").
				First(&user).Error; err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword(
				[]byte(user.Password),
				[]byte(payload.Password),
			); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &JwtPayload{
				UserID:   user.ID,
				Username: user.Username,
				IsAdmin:  user.IsAdmin,
			}, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if payload, ok := data.(*JwtPayload); ok {
				return jwt.MapClaims{
					"UserID":   payload.UserID,
					"Username": payload.Username,
					"IsAdmin":  payload.IsAdmin,
				}
			}
			fmt.Println("Empty Claim")
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			id, _ := claims["UserID"].(float64)
			username, _ := claims["Username"].(string)
			isAdmin, _ := claims["IsAdmin"].(bool)
			return &JwtPayload{
				UserID:   uint(id),
				Username: username,
				IsAdmin:  isAdmin,
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			url := c.Request.URL.String()
			adminOnly, err := regexp.MatchString("/admin", url)
			if err != nil {
				fmt.Println(err)
			}
			auth, ok := data.(*JwtPayload)
			if ok {
				if adminOnly {
					return auth.IsAdmin
				}
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})
}
