package controllers

import (
	"github.com/aiotrc/lanserver.sh/app/models"
	"github.com/revel/revel"
)

// Device controller controls system devices
type Device struct {
	*revel.Controller
}

// Create creates new device
func (c Device) Create() revel.Result {
	var d models.Device

	c.Params.BindJSON(&d)

	return c.RenderJSON(d)
}
