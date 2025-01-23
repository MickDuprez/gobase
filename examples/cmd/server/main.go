package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MickDuprez/gobase/core/app"
	"github.com/MickDuprez/gobase/core/config"
	"github.com/MickDuprez/gobase/core/utils"
	"github.com/MickDuprez/gobase/examples/features/about"
	"github.com/MickDuprez/gobase/examples/features/home"
	"github.com/MickDuprez/gobase/examples/features/users"
)

func main() {
	// for core framework development we need to set main directory to 'examples'
	// to maintain pathing for templates etc.
	os.Chdir("examples")

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
