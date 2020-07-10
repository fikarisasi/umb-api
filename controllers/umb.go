package controllers

import (
	"encoding/xml"
	"strconv"
	"umb_api/models"
	"strings"
	"regexp"
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

//  UmbController operations for Umb
type UmbController struct {
	beego.Controller
}

type locRequest struct {
	Tid    string `xml:"tid"`
	Msisdn string `xml:"msisdn"`
	Str    string `xml:"str"`
	V      string `xml:"v"`
	Action string `xml:"action"`
	Nodeid string `xml:"nodeid"`
}

type errorRes struct {
	Category   	string `json:"category"`
    Message 	error `json:"message"`
}

type Result struct {
	Value string `xml:"msisdn>CellID"`
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
	msisdn := c.Ctx.Input.Query("MSISDN")
	mid := c.Ctx.Input.Query("mid")
	sc := c.Ctx.Input.Query("sc")
	v, err := models.GetUmb(msisdn, mid, sc, mid, mid, mid, mid)
	// v, err := models.GetUmbById(id)
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
		Tid:    "123",
		Msisdn: strconv.FormatInt(id, 10),
		Str:    "EVENT",
		V:      "1",
		Action: "H%2780000000",
		Nodeid: "SDP",
	}
	req := httplib.Post(MapGatewayGenericUrl)
	req.XMLBody(body)
	str, err := req.String()
	if err != nil {
		mggErrRes := &errorRes {
			Category: "MapGatewayGeneric API Error",
			Message: err,
		}
		mggErrResJSON, _ := json.Marshal(mggErrRes)
    	beego.Error(string(mggErrResJSON))
		beego.Error(mggErrRes)
	}
	return str, nil
}

func (c *UmbController) GetMenu() {
	msisdn := c.Ctx.Input.Query("MSISDN")
	mid := c.Ctx.Input.Query("mid")
	sc := c.Ctx.Input.Query("sc")
	cell := c.Ctx.Input.Query("CELLID")
	regamtmn := c.Ctx.Input.Query("regamtmn")
	sms := c.Ctx.Input.Query("sms")
	userinput := c.Ctx.Input.Query("USERINPUT")
	sms = strings.Replace(sms, " ", "%20", -1)
	
	if mid == "OCODEPSS" && strings.Contains(sms, "%20") {
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		subString := strings.Split(sms, "%20")
		email := subString[5]
		
		if re.MatchString(email) {

		} else {
			mid = "EMAILPSS_FALSE"
		}
		beego.Info(email)
		beego.Info(mid)
	}

	if cell == "" {
		res := Result{}
		intMsisdn, _ := strconv.ParseInt(msisdn, 0, 64)
		data, _ := BackendData(intMsisdn)
		beego.Info(data)
		xml.Unmarshal([]byte(data), &res)
		beego.Info(res.Value)
		if res.Value == "" {
			cell = "999999"
		} else {
			cell = res.Value
		}
	}

	v, err := models.GetUmb(msisdn, mid, sc, cell, regamtmn, userinput, sms)
	if err != nil {
		c.Data["xml"] = err.Error()
	} else {
		c.Data["xml"] = v
	}
	c.ServeXML()
}
