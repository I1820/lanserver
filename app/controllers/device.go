package controllers

import (
	"net/http"

	"github.com/aiotrc/lanserver.sh/app"
	"github.com/aiotrc/lanserver.sh/app/models"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
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

	if err := app.DB.C("device").Insert(d); err != nil {
		c.Response.Status = http.StatusInternalServerError
		return revel.ErrorResult{
			Error: err,
		}
	}

	return c.RenderJSON(d)
}

// List lists all devices
func (c Device) List() revel.Result {
	var results []bson.M

	if err := app.DB.C("device").Find(bson.M{}).All(&results); err != nil {
		c.Response.Status = http.StatusInternalServerError
		return revel.ErrorResult{
			Error: err,
		}
	}

	return c.RenderJSON(results)
}
