package router

import (
	"context"
	"log"
	"time"

	"github.com/golang/glog"
	"github.com/openconfig/gnmic/pkg/api"
	"github.com/openconfig/gnmic/pkg/api/target"
	"github.com/openconfig/ygnmi/ygnmi"
)

type DeviceCfg struct {
	Address    string
	Username   string
	Password   string
	Port       string
	ID         string
	SkipVerify bool
	Timeout    time.Duration
}

type DeviceSession struct {
	ID     string
	Target *target.Target
	Client *ygnmi.Client
	State  *ConnectionState
}

func (m *DeviceManager) CreateNewDeviceSession(cfg DeviceCfg, ctx context.Context) (*DeviceSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// create target, connection to device firstly
	target, err := api.NewTarget(
		api.Address(cfg.Address+":"+cfg.Port),
		api.SkipVerify(cfg.SkipVerify),
		api.Timeout(cfg.Timeout),
		api.Username(cfg.Username),
		api.Password(cfg.Password),
	)
	if err != nil {
		return nil, err
	}

	// Create gnmi client
	err = target.CreateGNMIClient(ctx)
	if err != nil {
		return nil, err
	}

	// Create ygnmiClient by using gnmiClient
	client, err := createClient(&cfg, target)
	if err != nil {
		return nil, err
	}

	deviceSession := &DeviceSession{
		Target: target,
		ID:     cfg.ID,
		Client: client,
		State: &ConnectionState{
			Connection: "Connected",
			Telemetry:  "Running",
			Config:     "Idle",
		},
	}

	// adding this new session to device manager session list
	m.DeviceSessions[cfg.ID] = deviceSession

	return deviceSession, nil
}

func createClient(cfg *DeviceCfg, target *target.Target) (*ygnmi.Client, error) {
	client, err := ygnmi.NewClient(
		target.Client,
		ygnmi.WithRequestLogLevel(glog.Level(log.Default().Flags())),
		ygnmi.WithTarget(cfg.ID),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *DeviceSession) CloseClient() {
	s.Target.Close()
}
