package main

import (
	"log"

	"github.com/LeBaoTai/SDN-RD/internal/backend/handler"
	"github.com/LeBaoTai/SDN-RD/internal/backend/nats"
	"github.com/LeBaoTai/SDN-RD/internal/backend/router"
	"github.com/LeBaoTai/SDN-RD/internal/backend/service"
	"github.com/LeBaoTai/SDN-RD/internal/config"
)

func main() {
	cfg := config.LoadBEConfig()

	publisher, err := nats.NewPublisher(cfg.NatsURL)
	if err != nil {
		log.Fatalf("failed to connect NATS: %v", err)
	}

	intentService := service.NewIntentService(publisher)
	intentHandler := handler.NewIntentHandler(intentService)

	deps := router.RouterDeps{
		IntentHandler: intentHandler,
	}

	r := router.NewRouter(deps)
	if err := r.Run(cfg.BEAdd + ":" + cfg.BEPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
