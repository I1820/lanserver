package actions

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

// App creates new instance of Echo and configures it
func App(debug bool, db *mongo.Database) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Pre(middleware.RemoveTrailingSlash())

	if debug {
		app.Logger.SetLevel(log.DEBUG)
	}

	// validator
	app.Validator = &DefaultValidator{validator.New()}

	// routes
	app.GET("/about", AboutHandler)
	g := app.Group("/api")
	{
		dr := DevicesHandler{db}
		g.POST("/devices", dr.Create)
		g.GET("/devices", dr.List)
		g.GET("/devices/:device_id", dr.Show)
		g.PUT("/devices/:device_id", dr.Update)
		g.DELETE("/devices/:device_id", dr.Destroy)
		g.GET("/devices/:device_id/refresh", dr.Refresh)
	}

	return app
}
