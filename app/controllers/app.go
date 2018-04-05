package controllers

import (
	"github.com/revel/revel"
)

// App controller controls main functionality of application
type App struct {
	*revel.Controller
}

// Index method handles route "/"
func (c App) Index() revel.Result {
	return c.RenderText("18.20 is leaving us")
}
