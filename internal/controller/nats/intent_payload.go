package nats

import "encoding/json"

type IntentEnvelope struct {
	DeviceName string          `json:"device_name"`
	Type       string          `json:"type"`
	Data       json.RawMessage `json:"data"`
}

type IfcReq struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
	Duplex      string `json:"duplex"`
	SubIndex    uint32 `json:"sub-index"`
	Mtu         uint16 `json:"mtu"`
	Speed       int16  `json:"speed"`
	Mask        uint8  `json:"mask"`
	Enabled     bool   `json:"enable"`
}
