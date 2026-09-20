package controller

import (
	"log"
	"os"

	"go.yaml.in/yaml/v4"
)

type DeviceInfo struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type DeviceList struct {
	Devices []*DeviceInfo `yaml:"devices"`
}

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Start() {
	log.Println("Starting Controller......")

	// TODO: check queue

	// TODO: Check device list and establish connection
	checkDeviceList()

	log.Println("Controller is running")
}

func checkDeviceList() {
	devicesByte, err := os.ReadFile("../config/devices.yml")
	if err != nil {
		log.Fatalln("Cannot open DeviceFile: ", err)
	}
	var deviceList DeviceList
	err = yaml.Unmarshal(devicesByte, &deviceList)
	if err != nil {
		log.Fatalln("Error when parse yaml: ", err)
	}
	log.Println(deviceList)
}
