package controllers

import (
	"fmt"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    // "reflect"
    "strconv"
    models "umb_api/models"
)

type ArticleController struct {
    beego.Controller
}

func (c *ArticleController) Get() {
    orm.Debug = true
    o := orm.NewOrm()
    o.Using("default") 
    var maps []orm.Params
    // var Article models.Article
    _, err := o.QueryTable("Article").Values(&maps, "id", "name", "client", "url", "notes")
    if err == orm.ErrNoRows {
        fmt.Println("No records")
    } else if err == orm.ErrMissPK {
        fmt.Println("No Primary Key")
    } else {
        // for _, v := range maps {
        //     fmt.Println(v)
        // }
        c.Data["json"] = maps
        c.ServeJSON()
    } 
}

func (c *ArticleController) Read() {
    orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
    var id1 string= c.Ctx.Input.Param(":ids")
    id, err := strconv.Atoi(id1)
    user := models.Article{Id: id}
    err = o.Read(&user)
    if err == orm.ErrMultiRows {
        fmt.Printf("Returned Multi Rows Not One")
    }
    if err == orm.ErrNoRows {
        fmt.Printf("Not row found")
    }
    fmt.Println("id is: ", id)
    c.Data["json"] = user
    c.ServeJSON()
    // c.Data["m"] = user
    // c.TplName = "index.tpl"
}
