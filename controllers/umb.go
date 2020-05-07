package controllers

import (
	"strconv"
	"umb_api/models"

	"github.com/astaxie/beego"
)

//  UmbController operations for Umb
type UmbController struct {
	beego.Controller
}

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
	if err != nil {
		c.Data["xml"] = err.Error()
	} else {
		c.Data["xml"] = v
	}
	c.ServeXML()
}
