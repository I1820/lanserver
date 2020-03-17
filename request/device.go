package request

import validation "github.com/go-ozzo/ozzo-validation"

// Device represents device creation request
type Device struct {
	Name   string `json:"name"`
	DevEUI string `json:"devEUI"`
}

// Validate given device creation request
func (d Device) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.DevEUI, validation.Required),
	)
}
