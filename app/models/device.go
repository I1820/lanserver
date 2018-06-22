package models

import "net"

// Device that is connected by LAN
type Device struct {
	Name        string
	DevEUI      int64  // System wide identification
	IP          net.IP // network wide identification
	Token       string // Device JWT token
	Application string // Application ID
}

// DeviceProfile that is a profile for many connected device
type DeviceProfile struct {
	Name string
}
