package utils

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err:= godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or could not be loaded")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}