package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/LeBaoTai/myco-controller/internal/controller/builder"
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
		Username:   "admin",
		Password:   "NokiaSrl1!",
		SkipVerify: true,
		Port:       "57400",
		Timeout:    time.Second * 5,
	}

	targetConnection, err := gnmib.CreateTargetConnection(cfg, ctx)
	if err != nil {
		log.Fatalf("Cannot establish connection:%v", err)
	}

	log.Printf("Connected to: %v\n", targetConnection)
	mockData, err := os.Open("./mock-data/change.json")
	if err != nil {
		log.Printf("Cannot open the file: %v\n", err)
	}
	defer mockData.Close()

	byteValue, err := io.ReadAll(mockData)
	if err != nil {
		log.Printf("Cannot parse to byte: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌THE PROGRAM IS PANICED: %v", r)
		}
	}()

	var ifaceData builder.IfcReq
	err = json.Unmarshal(byteValue, &ifaceData)
	if err != nil {
		log.Printf("Error while Unmarshal, %v\n", err)
	}

	iface := builder.CreateInterface(&ifaceData)
	if err := iface.Validate(); err != nil {
		log.Printf("Invalid configuration: %v\n", err)
	} else {
		log.Println("Valid configuration")
	}

	result, err := builder.UpdateInterface(iface, targetConnection.Client, ctx)
	if err != nil {
		log.Println("Cannot update the interface: ", err)
	}
	log.Println("Result: ", result)
}
