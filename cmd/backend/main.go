package main

import (
	"log"

	backend "github.com/LeBaoTai/SDN-RD/internal/backend"
	"github.com/LeBaoTai/SDN-RD/internal/backend/handler"
	"github.com/LeBaoTai/SDN-RD/internal/backend/nats"
	"github.com/LeBaoTai/SDN-RD/internal/backend/service"
	"github.com/LeBaoTai/SDN-RD/internal/config"
)

func main() {
	cfg := config.Load()

	publisher, err := nats.NewPublisher(cfg.NatsURL)
	if err != nil {
		log.Fatalf("failed to connect NATS: %v", err)
	}

	intentService := service.NewIntentService(publisher)
	intentHandler := handler.NewIntentHandler(intentService)

	r := backend.NewRouter(intentHandler)
	if err := r.Run(cfg.BEAdd + ":" + cfg.BEPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
