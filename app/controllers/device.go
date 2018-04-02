package controllers

import (
	"net/http"

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

	if err := c.Params.BindJSON(&d); err != nil {
		c.Response.Status = http.StatusBadRequest
		return revel.ErrorResult{
			Error: err,
		}
	}

	return c.RenderJSON(d)
}
