package controllers

import (
	"context"
	"net/http"

	"github.com/aiotrc/lanserver.sh/app"
	"github.com/aiotrc/lanserver.sh/app/models"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
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

	res, err := app.DB.Collection("application").InsertOne(context.Background(), a)
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return revel.ErrorResult{
			Error: err,
		}
	}

	a.ID = res.InsertedID.(objectid.ObjectID).Hex()

	return c.RenderJSON(a)
}
