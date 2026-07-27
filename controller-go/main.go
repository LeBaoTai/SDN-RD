package main

import (
	"log"
	"time"

	"github.com/LeBaoTai/myco-controller/internal/controller/gnmib"
	"github.com/LeBaoTai/myco-controller/internal/controller/router"
)

type ConfigPath struct {
	Paths []string `yaml:"paths"`
}

func main() {
	client, err := gnmib.New(gnmib.Config{
		Address:    "192.100.100.101: 57400",
		Username:   "admin",
		Password:   "NokiaSrl1!",
		SkipVerify: true,
		Timeout:    5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	// newConfigFile, err := os.ReadFile("./mock-data/change.json")
	// if err != nil {
	// 	log.Printf("Cannot open newconfig file: %v", err)
	// }

	session := router.CreateNewSession(client)
	// session.HandleChageConfiguration(newConfigFile)

	currentConfig, err := session.GerCurrentConfiguration()
	if err != nil {
		log.Printf("Current config: %v\n", string(currentConfig))
	}
}
