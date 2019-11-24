package main

import (
	"fama-api/database"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		quitWithError(err)
	}
	err = run()
	if err != nil {
		quitWithError(err)
	}
	fmt.Println("All ok! Quiting")
}

func getOptions() string {
	return fmt.Sprintf(
		`host=%s port=%s dbname=%s user=%s password=%s`,
		os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"), os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"))
}

func run() error {
	options := getOptions()
	db, err := database.Connect(options)
	defer db.Close()
	if err != nil {
		return err
	}
	database.Migrate(db)

	fmt.Println("Connection OK!")
	return nil
}

func quitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
