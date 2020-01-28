package main

import (
	"fama-api/database/models"
	"fama-api/database"
	"github.com/jinzhu/gorm"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
	_ "net/http"
	_ "log"
	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		quitWithError(err)
	}
	err = run()
	if err != nil {
		quitWithError(err)
	}
	fmt.Println("All ok! Quiting")
}

func run() error {
	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		return err
	}

	fmt.Println("Connection OK!")
	err = listen(db)
	if err != nil {
		return err
	}
	return nil
}

func listen(db *gorm.DB) error {
	r := gin.Default()
	r.POST("/signup", func(c *gin.Context) {
		user := models.User{
			Username: c.PostForm("username"),
			Password: c.PostForm("password"),
			IsAdmin: false,
		}
		db.Create(&user)
		c.JSON(200, gin.H{"message": "Created new user!"})
	})
	r.POST("/login", func(c *gin.Context) {
		var user models.User
		db.Where(&models.User{Username: c.PostForm("username")}).Find(&user)
		c.JSON(200, gin.H{"user": user})
	})
	r.GET("/users", func(c *gin.Context) {
		var users []models.User
		db.Find(&users)
		c.JSON(200, gin.H{"users": users})
	})
	r.POST("/text", func(c *gin.Context) {
		text := models.Text{
			Body: c.PostForm("text"), 
		}
		db.Create(&text)
		c.JSON(200, gin.H{"message": "Created new text!"})
	})


	r.Run(":3000")
	return nil
}

func quitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
