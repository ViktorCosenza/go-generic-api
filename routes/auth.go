package routes

import (
	"fama-api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func createUser(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		type SignupPayload struct {
			Username string `form:"username" binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		var payload SignupPayload
		if err := c.Bind(&payload); err != nil {
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		user := models.User{
			Username: payload.Username,
			Password: string(hash),
		}
		if err := db.Create(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Username is taken"})
			return
		}
		db.Create(&models.Admin{UserID: user.ID}) // All users are admins for now TODO: REMOVE THIS
		c.JSON(http.StatusOK, gin.H{"Username": payload.Username, "Password": payload.Password, "Hash": string(hash)})
	}
}
