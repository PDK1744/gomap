package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type APIConfig struct {
	Addr string
}

func NewAPIConfig() *APIConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("ERROR: failed to load worker config: %v", err)
	}

	return &APIConfig{
		Addr: os.Getenv("API_PORT"),
	}
}
