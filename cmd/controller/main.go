package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/LeBaoTai/SDN-RD/internal/config"
	"github.com/LeBaoTai/SDN-RD/internal/controller/controller"
	"github.com/LeBaoTai/SDN-RD/internal/controller/nats"
	"github.com/LeBaoTai/SDN-RD/internal/shared/db"
	"github.com/LeBaoTai/SDN-RD/internal/shared/model"
	"github.com/LeBaoTai/SDN-RD/internal/shared/repo"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// load ENV file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Cannot load the env file %v", err)
	}
	ctlCfg := config.LoadCTLConfig()
	dbCfg := config.LoadDBConfig()

	// init DB connection
	db, err := db.NewConnection(*dbCfg)
	if err != nil {
		log.Fatalf("Cannot establish DB connection %v", err)
	}

	// Init repository
	repo := repo.NewRepo(db)

	// init controller
	controller := controller.NewController(repo)
	controller.Init(ctx, ctlCfg)

	sub, err := nats.NewSubscriber(ctlCfg.NatsURL)
	if err != nil {
		log.Fatalf("failed to connect NATS: %v", err)
	}
	defer sub.Close()
	err = sub.SubscribeQueue(ctlCfg.NatsSubject, ctlCfg.NatsQueueGroup, func(data []byte) {
		var payload model.IntentEnvelope
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

	log.Println("controller listening on subject:", ctlCfg.NatsSubject)
	select {}
}
