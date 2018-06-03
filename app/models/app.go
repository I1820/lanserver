package models

// Application is a collection of devices with specific purpose
type Application struct {
	Name        string `json:"name"`
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
}
