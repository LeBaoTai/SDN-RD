package controller

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/LeBaoTai/SDN-RD/internal/controller/controller/router"
	"go.yaml.in/yaml/v4"
	"google.golang.org/grpc/metadata"
)

type DeviceInfo struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

type DeviceList struct {
	Devices []*DeviceInfo `yaml:"devices"`
}

type DeviceCfg struct {
	Username   string        `yaml:"username"`
	Password   string        `yaml:"password"`
	SkipVerify bool          `yaml:"skipverify"`
	Timeout    time.Duration `yaml:"timeout"`
}

type Controller struct {
	DeviceSessions map[string]*router.DeviceSession
	mu             sync.Mutex
}

func NewController() *Controller {
	return &Controller{
		DeviceSessions: make(map[string]*router.DeviceSession),
	}
}

func (c *Controller) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Starting Controller......")

	// TODO: check queue

	// TODO: Check device list and establish connection
	deviceList := loadDeviceList()
	deviceCfg := loadDeviceConfig()

	authCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"username", deviceCfg.Username,
		"password", deviceCfg.Password,
	))

	c.establishConnection(deviceList, deviceCfg, authCtx)

	log.Println("Controller is running")
}

func (c *Controller) establishConnection(deviceList *DeviceList, devCfg *DeviceCfg, ctx context.Context) {
	log.Println("Establish connection to router....")
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, dev := range deviceList.Devices {
		devCfg := router.DeviceCfg{
			Address:    dev.Address,
			Username:   devCfg.Username,
			Password:   devCfg.Password,
			Port:       dev.Port,
			ID:         dev.Name,
			SkipVerify: devCfg.SkipVerify,
			Timeout:    devCfg.Timeout,
		}
		session, err := router.CreateNewDeviceSession(devCfg, ctx)
		if err != nil {
			log.Printf("Cannot establish connection with: %v\n", dev.Name)
			log.Println(err)
			continue
		}

		c.DeviceSessions[dev.Name] = session
	}
	log.Println("Completed establish connection to router....")
}

func loadDeviceList() *DeviceList {
	log.Println("Loading device list")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("Cannot get the current file information")
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
	log.Println("Completed loading device list")
	return &deviceList
}

func loadDeviceConfig() *DeviceCfg {
	log.Println("Loading Device Config")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("Cannot get the current file information")
	}
	currentDir := filepath.Dir(filename)
	deviceCfgPath := filepath.Join(currentDir, "../config/cfg.yml")
	deviceCfgByte, err := os.ReadFile(deviceCfgPath)
	if err != nil {
		log.Fatalln("Cannot open Device Config File: ", err)
	}
	var deviceCfg DeviceCfg
	err = yaml.Unmarshal(deviceCfgByte, &deviceCfg)
	if err != nil {
		log.Fatalln("Error when parse yaml: ", err)
	}
	log.Println("Completed loading Device Config")
	return &deviceCfg
}
