package app

import (
	"fmt"

	"github.com/revel/revel"
	"github.com/yosssi/gmq/mqtt/client"

	mgo "gopkg.in/mgo.v2"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

// DB is a connection to database
var DB *mgo.Database

// InitDB initiates a connection to a database
func InitDB() {
	url := revel.Config.StringDefault("db.url", "127.0.0.1")

	session, err := mgo.Dial(url)
	if err != nil {
		revel.AppLog.Errorf("DB connection error: %s", err)
		return
	}
	DB = session.DB("lanserver")
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// Device collection
	cp := DB.C("device")
	if err := cp.EnsureIndex(mgo.Index{
		Key:    []string{"deveui"},
		Unique: true,
	}); err != nil {
		revel.AppLog.Errorf("DB ensure index error: %s", err)
		return
	}

	revel.AppLog.Infof("DB Connected: %s", url)
}

// Secret key for signing jwt tokens
var Secret []byte

// InitSecret gets secret value from app.conf
func InitSecret() {
	Secret = []byte(revel.Config.StringDefault("app.secret", "shamin"))
}

// Mqtt client
var Mqtt *client.Client

// InitMQTT creates connection
func InitMQTT() {
	// Create an MQTT Client.
	Mqtt = client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Connect to the MQTT Server.
	if err := Mqtt.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  revel.Config.StringDefault("mqtt.url", "127.0.0.1:1883"),
		ClientID: []byte("lanserver.sh-client"),
	}); err != nil {
		revel.AppLog.Errorf("MQTT connection error: %s", err)
		return
	}
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	revel.OnAppStart(InitDB)
	revel.OnAppStart(InitSecret)
	revel.OnAppStart(InitMQTT)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
