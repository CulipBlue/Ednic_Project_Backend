package main

import (
	"log"

	_ "github.com/CulipBlue/backend_ednic/docs"
	"github.com/CulipBlue/backend_ednic/internal/config"
	"github.com/CulipBlue/backend_ednic/internal/database"
	"github.com/CulipBlue/backend_ednic/internal/httpapi"
)

// @title EDNIC Backend API
// @version 0.1.0
// @description API documentation for EDNIC backend.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use format: Bearer {token}
func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	router := httpapi.NewRouter(cfg, db)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
