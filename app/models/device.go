package models

// Device that is connected to LAN gateway
type Device struct {
	DevEUI  int64 // System wide identification
	DevAddr int   // network wide identification
}
