package model

import "encoding/json"

type IntentEnvelope struct {
	DeviceID   string          `json:"device_id" binding:"required"`
	DeviceName string          `json:"device_name" binding:"required"`
	Type       string          `json:"type" binding:"required"`
	Data       json.RawMessage `json:"data" binding:"required"`
}

type IntentResponse struct {
	Status  string
	Message string
}

type IfcReq struct {
	Name        string  `json:"name" binding:"required"`
	IP          string  `json:"ip" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Duplex      string  `json:"duplex" binding:"required"`
	Mtu         uint16  `json:"mtu" binding:"required"`
	Speed       int16   `json:"speed" binding:"required"`
	Mask        uint8   `json:"mask" binding:"required"`
	SubIndex    *uint32 `json:"sub-index" binding:"required"`
	Enabled     *bool   `json:"enable" binding:"required"`
}
