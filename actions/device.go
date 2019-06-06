package actions

import (
	"net/http"

	"github.com/I1820/lanserver/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type deviceReq struct {
	Name   string `json:"name" validate:"required"`
	DevEUI string `json:"devEUI" validate:"required"`
}

// DevicesHandler handles registered devices
type DevicesHandler struct {
	db *mongo.Database
}

// List gets all devices. This function is mapped to the path
// GET /devices
func (v DevicesHandler) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var results = make([]models.Device, 0)

	cur, err := v.db.Collection("devices").Find(ctx, bson.M{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for cur.Next(ctx) {
		var result models.Device

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// Show gets the data for one device. This function is mapped to
// the path GET /devices/{device_id}
func (v DevicesHandler) Show(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var d models.Device

	result := v.db.Collection("devices").FindOne(ctx, bson.M{
		"deveui": c.Param("device_id"),
	})

	if err := result.Decode(&d); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, d)
}

// Create adds a device to the DB. This function is mapped to the
// path POST /devices
func (v DevicesHandler) Create(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var d models.Device
	var rq deviceReq

	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	d.Name = rq.Name
	d.DevEUI = rq.DevEUI

	token, err := GenerateRandomString(32)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	d.Token = token

	if _, err := v.db.Collection("devices").InsertOne(ctx, d); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, d)
}

// Update changes a device in the DB. This function is mapped to
// the path PUT /devices/{device_id}
func (v DevicesHandler) Update(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var d models.Device
	var rq deviceReq

	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	d.Name = rq.Name

	res := v.db.Collection("devices").FindOneAndUpdate(ctx, bson.M{
		"deveui": c.Param("device_id"),
	}, bson.M{
		"$set": bson.M{
			"name": d.Name,
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if err := res.Decode(&d); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, d)
}

// Destroy deletes a device from the DB. This function is mapped
// to the path DELETE /devices/{device_id}
func (v DevicesHandler) Destroy(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var d models.Device

	result := v.db.Collection("devices").FindOneAndDelete(ctx, bson.M{
		"deveui": c.Param("device_id"),
	})

	if err := result.Decode(&d); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(200, d)
}

// Refresh creates new device token. This function is mapped to
// the path GET /devices/{device_id}/refresh
func (v DevicesHandler) Refresh(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	deveui := c.Param("device_id")

	token, err := GenerateRandomString(32)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if _, err := v.db.Collection("devices").UpdateOne(ctx, bson.M{
		"deveui": deveui,
	}, bson.M{
		"$set": bson.M{
			"token": token,
		},
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
