package config

import (
	"os"
)

type Config struct {
	NatsURL        string
	NatsSubject    string
	NatsQueueGroup string

	// Backend setting
	BEPort string
	BEAdd  string
}

func Load() *Config {
	return &Config{
		NatsURL:        getEnv("NATS_URL", "nats://127.0.0.1:4222"),
		NatsSubject:    getEnv("NATS_SUBJECT", "controller.intent"),
		NatsQueueGroup: getEnv("NATS_QUEUE_GROUP", "controller-workers"),
		BEPort:         getEnv("BEPort", "8080"),
		BEAdd:          getEnv("BEAdd", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
