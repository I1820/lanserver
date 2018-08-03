package models

// RxMessage represents recieved data from a device
type RxMessage struct {
	DevEUI string
	Data   []byte
}
