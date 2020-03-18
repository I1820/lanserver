package actions

import (
	"github.com/I1820/lanserver/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// App creates new instance of Echo and configures it
func App(debug bool, st store.Device) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Pre(middleware.RemoveTrailingSlash())

	if debug {
		app.Logger.SetLevel(log.DEBUG)
	}

	// routes
	app.GET("/about", AboutHandler)
	g := app.Group("/api")
	{
		dr := DevicesHandler{st}
		g.POST("/devices", dr.Create)
		g.GET("/devices", dr.List)
		g.GET("/devices/:device_id", dr.Show)
		g.PUT("/devices/:device_id", dr.Update)
		g.DELETE("/devices/:device_id", dr.Destroy)
		g.GET("/devices/:device_id/refresh", dr.Refresh)
	}

	return app
}
