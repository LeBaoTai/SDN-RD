package service

import (
	"log"

	"github.com/LeBaoTai/SDN-RD/internal/backend/model"
)

type IntentService struct {
	// publisher sẽ thêm lại sau khi cần NATS
}

func NewIntentService() *IntentService {
	return &IntentService{}
}

func (s *IntentService) HandleIntent(intent model.Intent) error {
	// TODO: gọi controller trực tiếp (in-process call hoặc HTTP) khi đã sẵn sàng
	log.Printf("received intent: %+v\n", intent)
	return nil
}
