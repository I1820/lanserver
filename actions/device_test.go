package actions

import "github.com/I1820/lanserver/models"

const (
	dName string = "test"
	dID   string = "0000000000000073"
	dIP   string = "192.168.73.10"
)

func (as *ActionSuite) Test_DevicesResource_Create_Show_Delete() {
	resc := as.JSON("/api/devices").Post(deviceReq{
		Name:   dName,
		DevEUI: dID,
		IP:     dIP,
	})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())

	var d models.Device

	ress := as.JSON("/api/devices/%s", "0000000000000073").Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())

	ress.Bind(&d)

	as.Equal(d.Name, dName)
	as.Equal(d.DevEUI, dID)
	as.Equal(d.IP.String(), dIP)

	resd := as.JSON("/api/devices/%s", "0000000000000073").Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}

func (as *ActionSuite) Test_DevicesResource_List() {
	var dl []models.Device

	res := as.JSON("/api/devices").Get()
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
	res.Bind(&dl)
}
