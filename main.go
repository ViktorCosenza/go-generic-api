package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"fama-api/server"
	_ "github.com/jinzhu/gorm/dialects/postgres"

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
