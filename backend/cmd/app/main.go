package main

import (
	"log"

	"github.com/almiluk/sipacks/config"
	"github.com/almiluk/sipacks/internal/app"
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
