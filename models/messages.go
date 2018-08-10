package models

// RxMessage represents recieved data from a device
// from lanserver into uplink
type RxMessage struct {
	DevEUI string
	Data   []byte
}

// TxMessage represents transmitted data to a device
// from downlink into downlink
type TxMessage struct {
	Confirmed bool
	FPort     int
	Data      []byte
}
