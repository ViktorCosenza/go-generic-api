package routes

import (
	"fama-api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

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

		var annotation models.Annotation
		if err := db.
			Model(&models.Annotation{}).
			Create(&models.Annotation{AssigmentID: payload.AssigmentID}).
			Find(&annotation).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var labels []interface{}
		for _, label := range payload.Labels {
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
	}
}
