package controllers

import (
	"fmt"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
	"strconv"
	"bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    models "umb_api/models"
)

type ConsumerController struct {
    beego.Controller
}

func (c *ConsumerController) Read() {
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
	c.ServeXML()
    // c.ServeJSON()
    // c.Data["m"] = user
    // c.TplName = "index.tpl"
}

func (c *ConsumerController) Test() {
	fmt.Println("Starting the application...")
    response, err := http.Get("https://httpbin.org/ip")
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }
    jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
    jsonValue, _ := json.Marshal(jsonData)
    response, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
        c.Data["json"] = data
        c.ServeJSON()
    }
	fmt.Println("Terminating the application...")
    // c.TplName = "views/index.tpl"
   
}