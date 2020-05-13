package routers

import (
    "github.com/astaxie/beego"
    "umb_api/controllers"
)

func init() {
    beego.Router("/read", &controllers.ArticleController{}, "*:Get")
    beego.Router("/read/:ids([0-9]+)", &controllers.ArticleController{}, "*:Read")
    beego.Router("/consumer", &controllers.ConsumerController{}, "*:Test")
    beego.Router("UMB/Menu", &controllers.UmbController{}, "*:GetMenu")
}