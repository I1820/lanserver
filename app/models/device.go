package models

import "net"

// Device that is connected by LAN
type Device struct {
	Name    string
	DevEUI  int64  // System wide identification
	DevAddr net.IP // network wide identification
	Token   string // Device JWT token
}

// DeviceProfile that is a profile for many connected device
type DeviceProfile struct {
}
