package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/LeBaoTai/myco-controller/internal/controller/builder"
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

	targetCfg := router.TargetCfg{
		Address:    "192.100.100.101",
		SkipVerify: true,
		Port:       "57400",
		Timeout:    time.Second * 5,
	}

	targetRouter, err := router.CreateNewTarget(targetCfg, authCtx)
	if err != nil {
		log.Println("Cannot create target: ", err)
	}

	clientRouter, err := router.CreateClient(targetRouter)
	if err != nil {
		log.Println("Cannot create client: ", err)
	}

	mockData, err := os.Open("./mock-data/change.json")
	if err != nil {
		log.Printf("Cannot open the file: %v\n", err)
	}
	defer mockData.Close()
	jsonData, err := io.ReadAll(mockData)
	if err != nil {
		log.Println("Cannot read data", err)
	}

	var ifaceReq builder.IfcReq
	err = json.Unmarshal(jsonData, &ifaceReq)
	if err != nil {
		log.Println("Cannot Unmarshal to struct", err)
	}
	log.Println(ifaceReq)
	ifaceUpdate, err := builder.CreateInterface(&ifaceReq)
	if err != nil {
		log.Println("Error when create configution", err)
	}

	result, err := builder.UpdateInterface(ifaceUpdate, clientRouter, authCtx)
	if err != nil {
		log.Println("Error when update interface: ", err)
	}
	log.Println("Update result: ", result)
}
