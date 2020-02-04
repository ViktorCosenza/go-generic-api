package routes

import (
	"archive/zip"
	"bytes"
	"fama-api/database/models"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

func getUsers(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var users []models.User
		db.Select("id, username, created_at, updated_at, deleted_at").Find(&users)
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func createTexts(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
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
	}
}

func createOntology(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		file, _, _ := c.Request.FormFile("file")
		json, _ := ioutil.ReadAll(file)
		db.Create(&models.JSONOntology{Value: string(json)})
		c.JSON(http.StatusOK, gin.H{"message": "Ontology saved."})
	}
}

func getAssigments(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
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
	}
}

func createAssigments(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
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
	}
}
