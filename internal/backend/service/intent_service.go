package service

import (
	"log"

	"github.com/LeBaoTai/SDN-RD/internal/backend/nats"
	"github.com/LeBaoTai/SDN-RD/internal/model"
)

type IntentService struct {
	Publisher *nats.Publisher
}

func NewIntentService(pub *nats.Publisher) *IntentService {
	return &IntentService{
		Publisher: pub,
	}
}

func (s *IntentService) HandleIntent(intent model.IntentEnvelope) error {
	log.Println(intent)
	return nil
}
