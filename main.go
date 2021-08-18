package main

import (
	"github.com/joho/godotenv"
	"github.com/muchlist/sagasql/app"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.RunApp()
}
