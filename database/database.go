package database

import (
	"github.com/jinzhu/gorm"
	"os"
)

// Connect to a database
func Connect() (*gorm.DB, error) {
	return gorm.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
}
