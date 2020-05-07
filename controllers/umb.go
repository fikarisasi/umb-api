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

// URLMapping ...
func (c *UmbController) URLMapping() {
	c.Mapping("GetOne", c.GetOne)
}

// GetOne ...
// @Title Get One
// @Description get Umb by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Umb
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UmbController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetUmbById(id)
	if err != nil {
		c.Data["xml"] = err.Error()
	} else {
		c.Data["xml"] = v
	}
	c.ServeXML()
}
