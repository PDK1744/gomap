package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	ConnStr string
}

func NewDBConfig() *DBConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("ERROR: could not load .env: %v", err)
	}
	mainConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	return &DBConfig{
		ConnStr: mainConnStr,
	}
}
