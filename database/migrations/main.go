package main

import (
	"fama-api/database"
	"fama-api/database/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"gopkg.in/gormigrate.v1"
	"log"
)

func main() {
	godotenv.Load()
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	db.LogMode(true)
	defer db.LogMode(false)

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "0001",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("users").Error
			},
		},
		{
			ID: "0002",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Text{}).Error
			},
		},
	})
	err = m.Migrate()
	if err != nil {
		log.Fatalf("Error in migration: %v", err)
	}
	log.Printf("Migration OK")
}
