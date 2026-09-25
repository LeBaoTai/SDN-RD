package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/LeBaoTai/SDN-RD/internal/controller/config"
	"github.com/LeBaoTai/SDN-RD/internal/controller/controller"
	"github.com/LeBaoTai/SDN-RD/internal/controller/nats"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// load ENV file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Cannot load the env file %v", err)
	}
	cfg := config.Load()

	controller := controller.NewController()
	controller.Init(ctx)

	sub, err := nats.NewSubscriber(cfg.NatsURL)
	if err != nil {
		log.Fatalf("failed to connect NATS: %v", err)
	}
	defer sub.Close()
	err = sub.SubscribeQueue(cfg.NatsSubject, cfg.NatsQueueGroup, func(data []byte) {
		var payload nats.IntentEnvelope
		if err := json.Unmarshal(data, &payload); err != nil {
			log.Printf("invalid payload: %v", err)
			return
		}

		if err := controller.ProcessIntent(ctx, &payload); err != nil {
			log.Printf("Process failed: %v", err)
			return
		}
		log.Printf("received intent: %+v\n", payload)
	})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	log.Println("controller listening on subject:", cfg.NatsSubject)
	select {}
}
