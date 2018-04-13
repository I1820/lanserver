package models

// Application is a collection of devices with specific purpose
type Application struct {
	Name        string
	ID          string `json:",omitempty"`
	Description string
}
