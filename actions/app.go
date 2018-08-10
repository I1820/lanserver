package actions

import (
	"context"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var db *mgo.Database
var mqttApplication paho.Client
var mqttNode paho.Client

var logger buffalo.Logger

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			// ADDR env
			// PORT env
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_lanserver_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))
		app.Use(func(next buffalo.Handler) buffalo.Handler {
			return func(c buffalo.Context) error {
				defer func() {
					c.Response().Header().Set("Content-Type", "application/json")
				}()

				return next(c)
			}
		})

		createMongodbClient()
		createMqttClient()

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		app.GET("/about", AboutHandler)
		g := app.Group("/api")
		{
			dr := DevicesResource{}
			g.Resource("/devices", dr)
			g.GET("/devices/{device_id}/refresh", dr.Refresh)
		}
	}

	return app
}

// createMongodbClient creates mongodb connection
func createMongodbClient() {
	url := envy.Get("DB_URL", "mongodb://127.0.0.1")
	client, err := mgo.NewClient(url)
	if err != nil {
		buffalo.NewLogger("fatal").Fatalf("DB new client error: %s", err)
	}
	if err := client.Connect(context.Background()); err != nil {
		buffalo.NewLogger("fatal").Fatalf("DB connection error: %s", err)
	}
	db = client.Database("lanserver")
}

// createMqttClient creates mqtt client and connect into broker
func createMqttClient() {
	// Application side mqtt
	{
		opts := paho.NewClientOptions()
		opts.AddBroker(envy.Get("APPLICATION_BROKER_URL", "tcp://127.0.0.1:1883"))
		mqttApplication = paho.NewClient(opts)
		if t := mqttApplication.Connect(); t.Wait() && t.Error() != nil {
			buffalo.NewLogger("fatal").Fatalf("MQTT session error: %s", t.Error())
		}
		mqttApplication.Subscribe("/device/+/tx", 0, Notification)
	}

	// Node side mqtt
	{
		opts := paho.NewClientOptions()
		opts.AddBroker(envy.Get("NODE_BROKER_URL", "tcp://127.0.0.1:1883"))
		mqttNode = paho.NewClient(opts)
		if t := mqttNode.Connect(); t.Wait() && t.Error() != nil {
			buffalo.NewLogger("fatal").Fatalf("MQTT session error: %s", t.Error())
		}
		mqttApplication.Subscribe("/log/+/send", 0, Log)
	}
}
