package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string, refrench string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Env is missing")
	}
	if vale := os.Getenv(key); vale != "" {
		return vale
	}
	return refrench
}
