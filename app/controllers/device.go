package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aiotrc/lanserver.sh/app"
	"github.com/aiotrc/lanserver.sh/app/models"
	jwt "github.com/dgrijalva/jwt-go"
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

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:   "lanserver.sh",
		Id:       strconv.FormatInt(d.DevEUI, 10),
		IssuedAt: time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(d.Token))
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return revel.ErrorResult{
			Error: err,
		}
	}

	d.Token = tokenString

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
