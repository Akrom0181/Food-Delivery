package main

import (
	"log"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
