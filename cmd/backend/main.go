package main

import (
	"log"

	backend "github.com/LeBaoTai/SDN-RD/internal/backend"
	"github.com/LeBaoTai/SDN-RD/internal/backend/config"
	"github.com/LeBaoTai/SDN-RD/internal/backend/handler"
	"github.com/LeBaoTai/SDN-RD/internal/backend/service"
)

func main() {
	cfg := config.Load()

	// publisher, err := nats.NewPublisher(cfg.NatsURL)
	// if err != nil {
	// 	log.Fatalf("failed to connect NATS: %v", err)
	// }

	intentService := service.NewIntentService()
	intentHandler := handler.NewIntentHandler(intentService)

	r := backend.NewRouter(intentHandler)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
