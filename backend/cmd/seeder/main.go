package main

import (
	"context"
	"log"

	"github.com/PDK1744/gomap/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer app.Close()
	log.Print("Starting Database seeding...")
	ctx := context.Background()
	if err := app.Branches.SeedMetadata(ctx, "stoneledge", "sl_branch"); err != nil {
		log.Fatalf("failed to seed metadata: %v", err)
	}

}
