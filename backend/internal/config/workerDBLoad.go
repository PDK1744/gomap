package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type WorkerDBConfig struct {
	ConnStr string
}

func NewWorkerDBConfig() *WorkerDBConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("ERROR: failed to load worker config: %v", err)
	}

	workerConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("WORKER_DB_USER"),
		os.Getenv("WORKER_DB_PASS"),
		os.Getenv("WORKER_DB_HOST"),
		os.Getenv("WORKER_DB_PORT"),
		os.Getenv("WORKER_DB_NAME"),
		os.Getenv("WORKER_DB_SSLMODE"),
	)

	return &WorkerDBConfig{
		ConnStr: workerConnStr,
	}
}
