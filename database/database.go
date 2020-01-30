package database

import (
	"os"

	"github.com/jinzhu/gorm"
)

// Connect to a database
func Connect() (*gorm.DB, error) {
	return gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
}
