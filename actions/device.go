package actions

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/I1820/lanserver/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
)

var devEUIRegexp *regexp.Regexp

type deviceReq struct {
	Name   string `json:"name"`
	DevEUI string `json:"devEUI"`
	IP     string `json:"ip"`
}

func init() {
	rg, err := regexp.Compile("[0-9a-fA-F]{16}")
	if err == nil {
		devEUIRegexp = rg
	}
}

// DevicesResource manages system devices
type DevicesResource struct {
	buffalo.Resource
}

// List gets all devices. This function is mapped to the path
// GET /devices
func (v DevicesResource) List(c buffalo.Context) error {
	var results = make([]models.Device, 0)

	cur, err := db.Collection("devices").Find(c, bson.NewDocument())
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for cur.Next(c) {
		var result models.Device

		if err := cur.Decode(&result); err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(200, r.JSON(results))
}

// Show gets the data for one device. This function is mapped to
// the path GET /devices/{device_id}
func (v DevicesResource) Show(c buffalo.Context) error {
	var d models.Device

	result := db.Collection("devices").FindOne(c, bson.NewDocument(
		bson.EC.String("deveui", c.Param("device_id")),
	))

	if err := result.Decode(&d); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(200, r.JSON(d))
}

// New renders the form for creating a new device.
// This function is mapped to the path GET /devices/new
func (v DevicesResource) New(c buffalo.Context) error {
	var rq deviceReq

	return c.Render(200, r.JSON(rq))
}

// Create adds a device to the DB. This function is mapped to the
// path POST /devices
func (v DevicesResource) Create(c buffalo.Context) error {
	var d models.Device
	var rq deviceReq

	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	d.Name = rq.Name
	if d.IP = net.ParseIP(rq.IP); d.IP == nil {
		return c.Error(http.StatusBadRequest, fmt.Errorf("Invalid ip address: %s", rq.IP))
	}
	if d.DevEUI = rq.DevEUI; !devEUIRegexp.MatchString(rq.DevEUI) {
		return c.Error(http.StatusBadRequest, fmt.Errorf("Invalid DevEUI: %s", rq.DevEUI))
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:   "lanserver.sh",
		Id:       d.DevEUI,
		IssuedAt: time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	var key []byte
	copy(key[:], envy.Get("SESSION_SECRET", ""))
	tokenString, err := token.SignedString(key)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	d.Token = tokenString

	if _, err := db.Collection("devices").InsertOne(c, d); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(200, r.JSON(d))
}

// Edit renders a edit formular for a device. This function is
// mapped to the path GET /deivces/{device_id}/edit
func (v DevicesResource) Edit(c buffalo.Context) error {
	var d models.Device

	result := db.Collection("devices").FindOne(c, bson.NewDocument(
		bson.EC.String("deveui", c.Param("device_id")),
	))

	if err := result.Decode(&d); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(200, r.JSON(d))
}

// Update changes a device in the DB. This function is mapped to
// the path PUT /devices/{device_id}
func (v DevicesResource) Update(c buffalo.Context) error {
	var d models.Device
	var rq deviceReq

	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	d.Name = rq.Name
	if d.IP = net.ParseIP(rq.IP); d.IP == nil {
		return c.Error(http.StatusBadRequest, fmt.Errorf("Invalid ip address: %s", rq.IP))
	}

	res := db.Collection("devices").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("deveui", c.Param("device_id")),
	), bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("name", d.Name),
			bson.EC.Interface("ip", d.IP),
		)),
	), findopt.ReturnDocument(mongoopt.After))

	if err := res.Decode(&d); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(200, r.JSON(d))
}

// Destroy deletes a device from the DB. This function is mapped
// to the path DELETE /devices/{device_id}
func (v DevicesResource) Destroy(c buffalo.Context) error {
	var d models.Device

	result := db.Collection("devices").FindOneAndDelete(c, bson.NewDocument(
		bson.EC.String("deveui", c.Param("device_id")),
	))

	if err := result.Decode(&d); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(200, r.JSON(d))
}

// Refresh creates new device token. This function is mapped to
// the path GET /devices/{device_id}/refresh
func (v DevicesResource) Refresh(c buffalo.Context) error {
	deveui := c.Param("device_id")

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:   "lanserver.sh",
		Id:       deveui,
		IssuedAt: time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	var key []byte
	copy(key[:], envy.Get("SESSION_SECRET", ""))
	tokenString, err := token.SignedString(key)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if _, err := db.Collection("devices").UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("deveui", deveui),
	), bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("token", tokenString),
		)),
	)); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(200, r.JSON(struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}))
}
