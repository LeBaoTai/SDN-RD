package controller

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/LeBaoTai/SDN-RD/internal/config"
	"github.com/LeBaoTai/SDN-RD/internal/controller/controller/router"
	"github.com/LeBaoTai/SDN-RD/internal/shared/model"
	"github.com/LeBaoTai/SDN-RD/internal/shared/repo"
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
	deviceSessions map[string]*router.DeviceSession
	repo           *repo.Repo
	mu             sync.Mutex
}

func NewController(rp *repo.Repo) *Controller {
	return &Controller{
		deviceSessions: make(map[string]*router.DeviceSession),
		repo:           rp,
	}
}

func (c *Controller) Init(ctx context.Context, cfg *config.CTLConfig) {
	log.Println("Starting Controller......")

	deviceList := loadDeviceList()

	authCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"username", cfg.SR_Username,
		"password", cfg.SR_Password,
	))

	c.establishConnection(deviceList, cfg, authCtx)

	log.Println("Controller is running")
}

func (c *Controller) establishConnection(deviceList *DeviceList, cfg *config.CTLConfig, ctx context.Context) {
	log.Println("Establish connection to router....")
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, dev := range deviceList.Devices {
		timeout, _ := time.ParseDuration(cfg.SR_Timeout)
		devCfg := router.DeviceCfg{
			Address:    dev.Address,
			Username:   cfg.SR_Username,
			Password:   cfg.SR_Password,
			Port:       dev.Port,
			Name:       dev.Name,
			SkipVerify: cfg.SR_Skipverify,
			Timeout:    timeout,
		}
		session, err := router.NewDeviceSession(devCfg, ctx)
		if err != nil {
			log.Printf("Cannot establish connection with: %v\n", dev.Name)
			log.Println(err)
			continue
		}

		c.deviceSessions[dev.Name] = session
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

func (c *Controller) LoadSession(s string) *router.DeviceSession {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := c.deviceSessions[s]
	return t
}

func (c *Controller) ProcessIntent(ctx context.Context, payload *model.IntentEnvelope) error {
	return nil
}
