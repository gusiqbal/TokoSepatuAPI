package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret []byte
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	secretString := os.Getenv("JWT_SECRET_KEY")

	if secretString == "" {
		log.Fatal("FATAL ERROR: JWT_SECRET belum di-set di file .env")
	}

	return &Config{
		JWTSecret: []byte(secretString),
	}
}
