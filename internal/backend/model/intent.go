package model

type Intent struct {
	DeviceName string            `json:"device_name" binding:"required"`
	Interface  string            `json:"interface" binding:"required"`
	IPAddress  string            `json:"ip_address" binding:"required"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type IntentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
