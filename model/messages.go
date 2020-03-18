package model

// RxMessage represents received data from a device
// from lanserver to uplink
type RxMessage struct {
	DevEUI string
	Data   []byte
}

// LogMessage represents raw data from a device
// from device to lanserver
type LogMessage struct {
	Data  []byte
	Token string
}

// TxMessage represents transmitted data to a device
// from downlink to downlink
type TxMessage struct {
	Confirmed bool
	FPort     int
	Data      []byte
}

// NotificationMessage represents raw data to a device
// from lanserver to device
type NotificationMessage struct {
	Data []byte
}
