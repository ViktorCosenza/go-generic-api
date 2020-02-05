package main

import (
	"fama-api/database"
	"fama-api/database/models"
	"flag"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"gopkg.in/gormigrate.v1"
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

func getMigrations(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "0001",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&models.User{}).Error; err != nil {
					return err
				}
				return tx.Model(&models.User{}).AddUniqueIndex("idx_username", "username").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("users").Error
			},
		},
		{
			ID: "0002",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&models.Admin{}).Error; err != nil {
					return err
				}

				if err := tx.Model(&models.Admin{}).
					AddUniqueIndex("idx_user_id", "user_id").Error; err != nil {
					return err
				}
				if err := tx.Model(models.Admin{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
					return err
				}

				password, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
				if err != nil {
					return err
				}

				var user models.User
				if err = tx.Table("users").
					Create(&models.User{
						Username: "Admin",
						Password: string(password)}).
					Scan(&user).Error; err != nil {
					return err
				}

				return tx.Table("admins").Create(&models.Admin{UserID: user.ID}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("admins").Error
			},
		},
		{
			ID: "0003",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&models.Text{}).Error; err != nil {
					return err
				}

				if err := tx.Model(models.Text{}).
					AddUniqueIndex("idx_name", "name").Error; err != nil {
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
		{
			ID: "0005",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&models.Assigment{}).Error; err != nil {
					return err
				}

				if err := tx.Model(&models.Assigment{}).
					AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
					return err
				}

				if err := tx.Model(&models.Assigment{}).
					AddUniqueIndex("idx_user_id_text_id", "user_id", "text_id").Error; err != nil {
					return err
				}

				return tx.Model(&models.Assigment{}).AddForeignKey("text_id", "texts(id)", "CASCADE", "CASCADE").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("assigments").Error
			},
		},
		{
			ID: "0007",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&models.Label{}).Error; err != nil {
					return err
				}
				return tx.Model(&models.Label{}).AddForeignKey("assigment_id", "assigments(id)", "CASCADE", "CASCADE").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("labels").Error
			},
		},
	})
}
