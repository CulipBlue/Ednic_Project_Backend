package main

import (
	"context"
	"log"
	"time"

	"github.com/CulipBlue/backend_ednic/internal/config"
	"github.com/CulipBlue/backend_ednic/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := database.RunMigrations(ctx, db, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("migrations applied successfully")
}
