package actions

import "github.com/aiotrc/lanserver/models"

func (as *ActionSuite) Test_DevicesResource_Create() {
	res := as.JSON("/api/devices").Post(deviceReq{
		Name:   "test",
		DevEUI: "0000000000000073",
		IP:     "192.168.73.10",
	})
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
}

func (as *ActionSuite) Test_DevicesResource_List() {
	var dl []models.Device

	res := as.JSON("/api/devices").Get()
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
	res.Bind(&dl)
}

func (as *ActionSuite) Test_DevicesResource_Show() {
	var d models.Device

	res := as.JSON("/api/devices/%s", "0000000000000073").Get()
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())

	res.Bind(&d)
}

func (as *ActionSuite) Test_DevicesResource_New() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_DevicesResource_Edit() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_DevicesResource_Update() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_DevicesResource_Destroy() {
	res := as.JSON("/api/devices/%s", "0000000000000073").Delete()
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
}
