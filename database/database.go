package database

import (
	"fama-api/database/models"
	"github.com/jinzhu/gorm"
)

// Connect to a database
func Connect(options string) (*gorm.DB, error) {
	return gorm.Open("postgres", options)
}

// Migrate performs all migrations
func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Text{},
		&models.Annotation{},
		&models.Label{},
	)
}
