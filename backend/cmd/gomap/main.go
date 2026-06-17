package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/PDK1744/gomap/internal/app"
)

// Database sync works but its also grabbin laptops because the laptop
// type is set to "pc" in the external database. Need to filter those out somehow
// And i only see PCs in the local assets table, no printers???
func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer app.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go app.StartAssetSyncWorker(ctx, 1*time.Hour)

	log.Printf("Starting server on %s", app.HttpServer.Addr)
	if err := app.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed to start: %v", err)
	}
}
