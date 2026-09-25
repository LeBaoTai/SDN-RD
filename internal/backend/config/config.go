package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	NatsURL string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	return Config{
		Port: getEnv("PORT", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
