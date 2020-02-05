package routes

import (
	"fama-api/database/models"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

func getAssigment(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		var assigment struct {
			ID   uint   `json:"assigment_id"`
			Body string `json:"text"`
			Name string `json:"title"`
		}
		if err := db.Table("assigments").
			Select("assigments.id, texts.body, texts.name").
			Joins("JOIN texts ON texts.id = assigments.text_id").
			Where(&models.Assigment{UserID: uint(claims["UserID"].(float64))}).
			Where("assigments.id NOT IN ?", db.Table("labels").Select("labels.assigment_id").Group("labels.assigment_id").SubQuery()).
			Take(&assigment).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "No assigments available for user", "assigment": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"assigment": assigment})
	}
}

func getOntology(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var ontology models.JSONOntology
		db.Find(&ontology)

		c.JSON(http.StatusOK, gin.H{"data": ontology.Value})
	}
}

func createAnnotation(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		type AnnotationPayload struct {
			Labels      []models.Label `json:"labels" binding:"required"`
			AssigmentID uint           `json:"assigment_id" binding:"required"`
		}

		var payload AnnotationPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var hasAnnotation uint
		if err := db.Table("labels").
			Where("labels.id = ?", payload.AssigmentID).
			Count(&hasAnnotation).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		if hasAnnotation > 0 {
			c.JSON(http.StatusConflict, gin.H{"err": "Assigment already fulfilled"})
		}
		var labels []interface{}
		for _, label := range payload.Labels {
			labels = append(labels, models.Label{
				First:       label.First,
				Second:      label.Second,
				Third:       label.Third,
				Fourth:      label.Fourth,
				Explicit:    label.Explicit,
				Start:       label.Start,
				End:         label.End,
				AssigmentID: payload.AssigmentID,
			})
		}

		if err := gormbulk.BulkInsert(db, labels, 3000); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": labels})
	}
}
