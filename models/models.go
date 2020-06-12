package models

import (
	"fmt"
	"os"

	"github.com/astaxie/beego/orm"

	_ "github.com/lib/pq"
)

type Article struct {
	Id     int    `orm:"pk"`
	Name   string `orm:"name,text,name:" valid:"MinSize(5);MaxSize(20)"`
	Client string `orm:"client,text,client:"`
	Url    string `orm:"url,text,url:"`
	Notes  string `orm:"url,text,notes:"`
}

type UmbHeader struct {
	MenuId     string `orm:"column(menu_id);pk"`
	MenuHeader string `orm:"column(menu_header)"`
}

func (header *UmbHeader) TableName() string {
	return "service_dyn_umb_header"
}

type UmbMenu struct {
	MenuId         string `orm:"column(menu_id);pk"`
	MenuDetailItem string `orm:"column(menu_detail_item)"`
}

func (menu *UmbMenu) TableName() string {
	return "service_dyn_umb_menu"
}

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase(os.Getenv("postgresql_schema"), "postgres", "user="+os.Getenv("postgresql_database_user")+" password="+os.Getenv("postgresql_database_password")+" host="+os.Getenv("postgresql_host")+" port="+"5432"+" dbname="+os.Getenv("postgresql_database_name")+" sslmode=disable")
	orm.RegisterDataBase(os.Getenv("postgresql_database_name"), "postgres", "user="+os.Getenv("postgresql_database_user")+" password="+os.Getenv("postgresql_database_password")+" host="+os.Getenv("postgresql_host")+" port="+"5432"+" dbname="+os.Getenv("postgresql_database_name")+" sslmode=disable")
	orm.RegisterModel(new(Article), new(UmbHeader), new(UmbMenu))
	fmt.Println("------------Setting schema--------------------")
	fmt.Println("user=" + os.Getenv("postgresql_database_user") + " password=" + os.Getenv("postgresql_database_password") + " host=" + os.Getenv("postgresql_host") + " port=" + "5432" + " dbname=" + os.Getenv("postgresql_database_name") + " sslmode=disable")
	//设置scheme
	o := orm.NewOrm()
	o.Using(os.Getenv("postgresql_database_name")) // Using public, you can use other database
	_, e := o.Raw("set search_path to ssp").Exec()
	if e != nil {
		panic(e)
	}

	orm.RunSyncdb(os.Getenv("postgresql_schema"), false, true)
	fmt.Println("------------Setting schema Completed--------------------")
}
