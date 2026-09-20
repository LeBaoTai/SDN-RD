package controller

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

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
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("Không thể lấy thông tin file hiện tại")
	}

	currentDir := filepath.Dir(filename)
	devicesPath := filepath.Join(currentDir, "../config/devices.yml")
	devicesByte, err := os.ReadFile(devicesPath)
	if err != nil {
		log.Fatalln("Cannot open Device File: ", err)
	}
	var deviceList DeviceList
	err = yaml.Unmarshal(devicesByte, &deviceList)
	if err != nil {
		log.Fatalln("Error when parse yaml: ", err)
	}
	for _, val := range deviceList.Devices {
		log.Println(val.Address)
	}
}
