package config

import (
	"github.com/joho/godotenv"
	"log"
)

func Init() {
	err := godotenv.Load("internal/config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Successfully loaded .env file")
}
