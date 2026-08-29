package main

import (
	"context"
	"log"
	"time"

	"github.com/LeBaoTai/myco-controller/internal/controller/gnmib"
)

type ConfigPath struct {
	Paths []string `yaml:"paths"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := gnmib.Cfg{
		Address:    "192.100.100.101",
		Password:   "NokiaSrl1!",
		SkipVerify: true,
		Timeout:    time.Minute * 5,
	}

	connection, err := gnmib.CreateTargetConnection(cfg, ctx)
	if err != nil {
		log.Printf("Cannot establish connection:%v", err)
	}
}
