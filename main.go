package main

import (
	"fama-api/server"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	err := run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("All ok! Quiting")
	}
}

func run() error {
	s, err := server.Start()
	if err != nil {
		return err
	}
	s.Listen(":8080")
	return nil
}
