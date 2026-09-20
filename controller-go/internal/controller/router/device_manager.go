package router

import (
	"sync"
)

type DeviceManager struct {
	DeviceSessions map[string]*DeviceSession
	mu             sync.Mutex
}

func CreateDeviceManager() *DeviceManager {
	return &DeviceManager{
		DeviceSessions: make(map[string]*DeviceSession),
	}
}
