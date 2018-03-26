package models

// Device that is connected to LAN gateway
type Device struct {
	Name    string
	DevEUI  int64 // System wide identification
	DevAddr int   // network wide identification
}
