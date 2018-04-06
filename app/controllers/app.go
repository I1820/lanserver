package controllers

import (
	"net/http"

	"github.com/aiotrc/lanserver.sh/app/models"
	"github.com/revel/revel"
)

// App controller controls main functionality of application
type App struct {
	*revel.Controller
}

// Create creates new application
func (c *App) Create() revel.Result {
	var a models.Application

	if err := c.Params.BindJSON(&a); err != nil {
		c.Response.Status = http.StatusBadRequest
		return revel.ErrorResult{
			Error: err,
		}
	}

	return c.RenderJSON(a)
}
