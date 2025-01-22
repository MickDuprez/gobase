package main

import (
	"log"
	"net/http"

	"github.com/MickDuprez/gobase/internal/core/app"
	"github.com/MickDuprez/gobase/internal/core/config"
	"github.com/MickDuprez/gobase/internal/core/utils"
	"github.com/MickDuprez/gobase/internal/features/about"
	"github.com/MickDuprez/gobase/internal/features/home"
	"github.com/MickDuprez/gobase/internal/features/users"
)

func main() {
	// Load .env file first
	if err := utils.LoadEnvFile(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Create app config based on environment setting in .env
	cfg := config.NewAppConfig()

	// Initialize app
	app, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Register features
	if err := app.RegisterFeature(home.New()); err != nil {
		log.Fatal(err)
	}
	if err := app.RegisterFeature(about.New()); err != nil {
		log.Fatal(err)
	}
	if err := app.RegisterFeature(users.New()); err != nil {
		log.Fatal(err)
	}

	// Start server
	log.Printf("Server starting on %s", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, app); err != nil {
		log.Fatal(err)
	}
}
