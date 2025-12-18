package main

import (
	"log"

	"kafka-consumer-challenge/app/events"
	"kafka-consumer-challenge/internal/config"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Start Application
	event.Start(cfg)
}
