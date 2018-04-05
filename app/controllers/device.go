package controllers

import (
	"fmt"
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
	tokenString, err := token.SignedString(app.Secret)
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

// Remove removes given device
func (c Device) Remove() revel.Result {
	deveui, err := strconv.ParseInt(c.Params.Route.Get("id"), 10, 64)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return revel.ErrorResult{
			Error: err,
		}
	}

	if err := app.DB.C("device").Remove(bson.M{
		"deveui": deveui,
	}); err != nil {
		c.Response.Status = http.StatusInternalServerError
		return revel.ErrorResult{
			Error: err,
		}
	}

	return c.RenderText(strconv.FormatInt(deveui, 10))
}

// Push stores device given data
func (c Device) Push() revel.Result {
	id := c.Params.Route.Get("id")

	var d struct {
		Token string
		Data  []byte
	}

	if err := c.Params.BindJSON(&d); err != nil {
		c.Response.Status = http.StatusBadRequest
		return revel.ErrorResult{
			Error: err,
		}
	}

	token, err := jwt.ParseWithClaims(d.Token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		c := token.Claims.(*jwt.StandardClaims)

		if !c.VerifyIssuer("lanserver.sh", true) {
			return nil, fmt.Errorf("Unexpected issuer %v", c.Issuer)
		}
		if c.Id != id {
			return nil, fmt.Errorf("Mismatched identifier %s != %s", c.Id, id)
		}
		return app.Secret, nil
	})
	if err != nil {
		c.Response.Status = http.StatusForbidden
		return revel.ErrorResult{
			Error: err,
		}
	}

	return c.RenderJSON(token.Claims)
}
