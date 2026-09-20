package main

import (
	"context"
	"log"
	"time"

	"github.com/LeBaoTai/myco-controller/internal/controller/router"
	"google.golang.org/grpc/metadata"
)

type ConfigPath struct {
	Paths []string `yaml:"paths"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	authCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"username", "admin",
		"password", "NokiaSrl1!",
	))

	deviceCfg := router.DeviceCfg{
		Address:    "192.100.100.101",
		SkipVerify: true,
		Port:       "57400",
		Timeout:    time.Second * 5,
		Username:   "admin",
		Password:   "NokiaSrl1!",
		ID:         "router-core",
	}

	deviceManager := router.CreateDeviceManager()
	newDeviceSession, err := deviceManager.CreateNewDeviceSession(deviceCfg, authCtx)
	if err != nil {
		log.Println("Cannot create a new sesion: ", err)
	}
	log.Println("New session: ", newDeviceSession.State.CheckState())
}
