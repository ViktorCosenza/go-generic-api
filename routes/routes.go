package routes

import(
	"fama-api/database/models"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"fmt"
	"archive/zip"
	"bytes"
)

// Start returns a router with the routes
func Start(db *gorm.DB) *gin.Engine {
	db.LogMode(true)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	r.POST("/signup", func(c *gin.Context) {
		type SignupPayload struct {
			Username string `form:"username" binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		if err := c.Bind(&SignupPayload{}); err != nil {
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		user := models.User{
			Username: c.PostForm("username"),
			Password: string(hash),
		}
		db.Create(&user)
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.POST("/login", func(c *gin.Context) {
		type LoginPayload struct {
			Username string `form:"username" binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		if err := c.Bind(&LoginPayload{}); err != nil {
			return
		}
		var user models.User
		db.Where(&models.User{Username: c.PostForm("username")}).First(&user)
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.GET("/users", func(c *gin.Context) {
		var users []models.User
		db.Find(&users)
		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	r.POST("/text", func(c *gin.Context) {
		var response []string
		zipFile, header, _ := c.Request.FormFile("file")
		read, _ := zip.NewReader(zipFile, header.Size)
		for _, file := range read.File {
			fileread, err := file.Open()
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			}
			defer fileread.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(fileread)
			response = append(response, buf.String())
		}
		c.JSON(http.StatusOK, gin.H{"message": response})
	})

	r.POST("/ontology", func(c *gin.Context) {
		file, _, _ := c.Request.FormFile("file")
		json, _ := ioutil.ReadAll(file)
		db.Create(&models.JSONOntology{Value: string(json)})
		c.JSON(http.StatusOK, gin.H{"file": string(json)})
	})

	r.GET("/ontology", func(c *gin.Context) {
		var ontology models.JSONOntology
		db.Find(&ontology)

		c.JSON(http.StatusOK, gin.H{"data": ontology.Value})
	})

	return r
}