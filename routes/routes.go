package routes

import (
	"archive/zip"
	"bytes"
	"fama-api/database/models"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
	"golang.org/x/crypto/bcrypt"
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
			return
		}
		user := models.User{
			Username: c.PostForm("username"),
			Password: string(hash),
		}
		if err := db.Create(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Username is taken"})
			return
		}
		db.Create(&models.Admin{UserID: user.ID}) // All users are admins for now TODO: REMOVE THIS
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
	}) // Get all users, just for dev purpose TODO: DELETE THIS WHEN DONE//

	// Add Texts via Zip file upload
	r.POST("/text", func(c *gin.Context) {
		var texts []interface{}
		zipFile, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		read, err := zip.NewReader(zipFile, header.Size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		for _, file := range read.File {
			fileread, err := file.Open()
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			}
			defer fileread.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(fileread)
			texts = append(texts, &models.Text{Body: buf.String(), Name: file.Name, AdminID: 1})
		}
		err = gormbulk.BulkInsert(db, texts, 1000)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "All files saved."})
	})

	r.POST("/ontology", func(c *gin.Context) {
		file, _, _ := c.Request.FormFile("file")
		json, _ := ioutil.ReadAll(file)
		db.Create(&models.JSONOntology{Value: string(json)})
		c.JSON(http.StatusOK, gin.H{"message": "Ontology saved."})
	})

	r.GET("/ontology", func(c *gin.Context) {
		var ontology models.JSONOntology
		db.Find(&ontology)

		c.JSON(http.StatusOK, gin.H{"data": ontology.Value})
	})

	r.GET("/assigment", func(c *gin.Context) { // Get User assign counts
		type AssigmentsCounts struct {
			Username string
			Count    uint
		}
		var assigmentCounts []AssigmentsCounts
		db.
			Table("users").
			Select("users.username, COUNT(assigments.user_id) as count").
			Joins("JOIN assigments ON assigments.user_id = users.id").
			Group("assigments.user_id, users.username").
			Find(&assigmentCounts)

		c.JSON(http.StatusOK, gin.H{"data": assigmentCounts})
	})

	r.POST("/assigment", func(c *gin.Context) {
		type AssigmentPayload struct {
			UserIDs      []uint `json:"user_ids" binding:"required"`
			TextQuantity uint   `json:"text_quantity" binding:"required"`
		}
		var payload AssigmentPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error()})
			return
		}

		var users []models.User
		var unassignedTexts []models.Text
		db.Where(payload.UserIDs).
			Select("id").
			Find(&users)
		db.Limit(payload.TextQuantity).
			Select("id").
			Where("texts.id NOT IN ?", db.Model(&models.Assigment{}).Select("text_id").Group("text_id").SubQuery()).
			Find(&unassignedTexts)
		var assigments []interface{}
		for _, user := range users {
			for _, text := range unassignedTexts {
				a := models.Assigment{UserID: user.ID, TextID: text.ID}
				assigments = append(assigments, a)
			}
		}

		err := gormbulk.BulkInsert(db, assigments, 2000)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Successfully assigned"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"assigments": assigments})
	})

	// Submit annotation for a given text
	r.POST("/annotation", func(c *gin.Context) {
		type AnnotationPayload struct {
			Annotation struct {
				Labels      []models.Label `json:"labels" binding:"required"`
				AssigmentID uint           `json:"assigment" binding:"required"`
			} `json:"annotation" binding:"required"`
		}

		var payload AnnotationPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		var annotation models.Annotation
		db.
			Model(&models.Annotation{}).
			Create(&models.Annotation{AssigmentID: payload.Annotation.AssigmentID}).
			Find(&annotation)

		var labels []interface{}
		for _, label := range payload.Annotation.Labels {
			labels = append(labels, models.Label{
				First:        label.First,
				Second:       label.Second,
				Third:        label.Third,
				Fourth:       label.Fourth,
				Explicit:     label.Explicit,
				Start:        label.Start,
				End:          label.End,
				AnnotationID: annotation.ID,
			})
		}

		if err := gormbulk.BulkInsert(db, labels, 3000); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": labels})
	})

	return r
}
