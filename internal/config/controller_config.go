package config

import "strconv"

type CTLConfig struct {
	NatsURL        string
	NatsSubject    string
	NatsQueueGroup string

	// Backend setting
	SR_Username   string
	SR_Password   string
	SR_Skipverify bool
	SR_Timeout    string
}

func LoadCTLConfig() *CTLConfig {
	skipVerify, _ := strconv.ParseBool(getEnv("SR_SKIPVERIFY", "true"))
	return &CTLConfig{
		NatsURL:        getEnv("NATS_URL", "nats://127.0.0.1:4222"),
		NatsSubject:    getEnv("NATS_SUBJECT", "controller.intent"),
		NatsQueueGroup: getEnv("NATS_QUEUE_GROUP", "controller-workers"),
		SR_Username:    getEnv("SR_USERNAME", "admin"),
		SR_Password:    getEnv("SR_PASSWORD", "NokiaSrl1!"),
		SR_Skipverify:  skipVerify,
	}
}
