package actions

import (
	"net/http"

	"github.com/I1820/lanserver/model"
	"github.com/I1820/lanserver/request"
	"github.com/I1820/lanserver/store"
	"github.com/labstack/echo/v4"
)

const (
	// TokenLength is device's token length
	TokenLength = 32
)

// DevicesHandler handles registered devices
type DevicesHandler struct {
	Store store.Device
}

// List gets all devices. This function is mapped to the path
// GET /devices
func (v DevicesHandler) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	results, err := v.Store.Get(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// Show gets the data for one device. This function is mapped to
// the path GET /devices/{device_id}
func (v DevicesHandler) Show(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	d, err := v.Store.Show(ctx, c.Param("device_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, d)
}

// Create adds a device to the DB. This function is mapped to the
// path POST /devices
func (v DevicesHandler) Create(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var rq request.Device

	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := GenerateRandomString(TokenLength)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	d := model.Device{
		Name:   rq.Name,
		DevEUI: rq.DevEUI,
		Token:  token,
	}

	if err := v.Store.Insert(ctx, d); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, d)
}

// Update changes a device in the DB. This function is mapped to
// the path PUT /devices/{device_id}
func (v DevicesHandler) Update(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var rq request.Device

	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	d, err := v.Store.Update(ctx, c.Param("device_id"), "name", rq.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, d)
}

// Destroy deletes a device from the DB. This function is mapped
// to the path DELETE /devices/{device_id}
func (v DevicesHandler) Destroy(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	d, err := v.Store.Destroy(ctx, c.Param("device_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, d)
}

// Refresh creates new device token. This function is mapped to
// the path GET /devices/{device_id}/refresh
func (v DevicesHandler) Refresh(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	deveui := c.Param("device_id")

	token, err := GenerateRandomString(TokenLength)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if _, err := v.Store.Update(ctx, deveui, "token", token); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
