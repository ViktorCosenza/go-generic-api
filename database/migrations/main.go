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
	"flag"
)

func main() {
	reset := flag.Bool("reset", false, "Reset all database before running migration")
	flag.Parse()
	godotenv.Load()
	db, err := database.Connect()
	if err != nil {
		fmt.Println("Could not connect: ", err)
		return
	}
	db.LogMode(true)

	m := getMigrations(db)

	if *reset {
		log.Println("Reseting database...")
		err = m.RollbackTo("0001")
		if err != nil {
			log.Fatalf("Migration error: %v", err)
		}
	}

	log.Println("Migrating...")
	err = m.Migrate()

	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	log.Println("Migration OK")
}

func getMigrations(db *gorm.DB) *gormigrate.Gormigrate{
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "0001",
			Migrate: func(tx *gorm.DB) error {				
				err := tx.AutoMigrate(&models.User{}).Error
				return err
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("users").Error
			},
		},
		{
			ID:"0002",
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&models.Admin{}).Error
				if err != nil {
					return err
				}
				return tx.Model(models.Admin{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("admins").Error
			},
		},
		{
			ID: "0003",
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&models.Text{}).Error
				if err != nil {
					return err
				}

				return tx.Model(models.Text{}).AddForeignKey("admin_id", "admins(id)", "CASCADE", "CASCADE").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("texts").Error
			},
		},
		{
			ID: "0004",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.JSONOntology{}).Error
			},
			Rollback: func(tx *gorm.DB) error { 
				return tx.DropTableIfExists("json_ontologies").Error
			},
		},
	})
}
