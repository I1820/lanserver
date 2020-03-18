package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/I1820/lanserver/model"
	"github.com/I1820/lanserver/request"
	"github.com/labstack/echo/v4"
)

const (
	dName string = "raha"
	dID   string = "0000000000000073"
)

func (suite *LSTestSuite) devicesHandlerCreate() {
	data, err := json.Marshal(request.Device{
		Name:   dName,
		DevEUI: dID,
	})
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"POST",
		"/api/devices",
		bytes.NewReader(data),
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *LSTestSuite) devicesHandlerShow() {
	var d model.Device

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("/api/devices/%s", dID),
		nil,
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &d))

	suite.Equal(d.Name, dName)
	suite.Equal(d.DevEUI, dID)
}

func (suite *LSTestSuite) devicesHandlerUpdate() {
	var d model.Device

	data, err := json.Marshal(request.Device{
		Name: "elahe",
	})
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/devices/%s", dID),
		bytes.NewReader(data),
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &d))

	suite.Equal(d.Name, "elahe")
	suite.Equal(d.DevEUI, dID)
}

func (suite *LSTestSuite) devicesHandlerDelete() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/devices/%s", dID),
		nil,
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *LSTestSuite) devicesHandlerList() {
	var dl []model.Device

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"GET",
		"/api/devices",
		nil,
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &dl))
}

func (suite *LSTestSuite) TestDevicesHandler() {
	suite.Run("create", suite.devicesHandlerCreate)
	suite.Run("show", suite.devicesHandlerShow)
	suite.Run("update", suite.devicesHandlerUpdate)
	suite.Run("delete", suite.devicesHandlerDelete)
	suite.Run("delete", suite.devicesHandlerList)
}
