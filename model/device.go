package models

// Device that is connected by LAN
type Device struct {
	Name   string `json:"name"`
	DevEUI string `json:"devEUI"`
	Token  string `json:"token"`
}

// Devices is not required by pop and may be deleted
type Devices []Device
