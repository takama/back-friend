package main

import (
	"log"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/service"
)

func main() {
	cfg := config.New()

	// Load configuration from environment vars
	if err := cfg.LoadFromEnv(); err != nil {
		log.Fatal("Error loading of configuration from environment: ", err)
	}

	if err := service.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
