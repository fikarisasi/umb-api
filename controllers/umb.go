package controllers

import (
	"strconv"
	"umb_api/models"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego"
	"encoding/xml"
)

//  UmbController operations for Umb
type UmbController struct {
	beego.Controller
}

type locRequest struct {
	Tid  	string `xml:"tid"`
	Msisdn  string `xml:"msisdn"`
	Str 	string `xml:"str"`
	V 		string `xml:"v"`
	Action	string `xml:"action"`
	Nodeid 	string `xml:"nodeid"`
}

type Result struct {
	Value	string	`xml:"msisdn>CellID"`
}

var MapGatewayGenericUrl = "http://10.34.234.180:8023/mapgw_generic/request_handler"

// GetOne ...
// @Title Get One
// @Description get Umb by id
// @Param	MSISDN		query 	string	true		"MSISDN"
// @Param	mid		query 	string	true		"mid"
// @Param	sc		query 	string	true		"sc"
// @Success 200 {object} models.Umb
// @Failure 403 MSISDN is empty
// @router / [get]
func (c *UmbController) GetOne() {
	idStr := c.Ctx.Input.Query("MSISDN")
	beego.Info(idStr)
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetUmbById(id)
	data, err := BackendData(id)
	if err != nil {
		c.Data["xml"] = err.Error()
	} else {
		c.Data["xml"] = v
	}
	res := Result{}
	xml.Unmarshal([]byte(data), &res)
	beego.Info(res.Value)
	c.ServeXML()
}

// Get data from backend API MapGatewayGeneric
func BackendData(id int64) (str string, err2 error) {
	body := &locRequest{
		Tid: "123", 
		Msisdn: strconv.FormatInt(id, 10), 
		Str: "EVENT", 
		V: "1", 
		Action: "H%2780000000", 
		Nodeid: "SDP",
	}
	req := httplib.Post(MapGatewayGenericUrl)
	req.XMLBody(body)
	str, err := req.String()
	if err != nil {
		beego.Info(err)
	}
	return str, nil
}