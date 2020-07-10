package models

import (
	"fmt"

	"github.com/astaxie/beego"
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
	ItemNumber 	   string `orm:"column(item_number)"`
	MenuNextId 	   string `orm:"column(menu_next_id)"`
	RegAmount 	   string `orm:"column(reg_amount)"`
	Unit 		   string `orm:"column(unit)"`
	Formula 	   string `orm:"column(formula)"`
	Keyword 	   string `orm:"column(keyword)"`
}

func (menu *UmbMenu) TableName() string {
	return "service_dyn_umb_menu"
}

type PromoCodeSSP struct {
	MenuId         string `orm:"column(menu_id)"`
	PromoCode      string `orm:"column(promo_code);pk"`
}

func (promoCodeDB *PromoCodeSSP) TableName() string {
	return "service_dyn_umb_prmcode"
}

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase(os.Getenv("postgresql_schema"), "postgres", "user="+os.Getenv("postgresql_database_user")+" password="+os.Getenv("postgresql_database_password")+" host="+os.Getenv("postgresql_host")+" port="+"5432"+" dbname="+os.Getenv("postgresql_database_name")+" sslmode=disable")
	orm.RegisterDataBase(os.Getenv("postgresql_database_name"), "postgres", "user="+os.Getenv("postgresql_database_user")+" password="+os.Getenv("postgresql_database_password")+" host="+os.Getenv("postgresql_host")+" port="+"5432"+" dbname="+os.Getenv("postgresql_database_name")+" sslmode=disable")
	orm.RegisterModel(new(Article), new(UmbHeader), new(UmbMenu), new(PromoCodeSSP))
	fmt.Println("------------Setting schema--------------------")
	//设置scheme
	o := orm.NewOrm()
	o.Using("umbdynamicdb") // Using public, you can use other database
	_, e := o.Raw("set search_path to ssp").Exec()
	if e != nil {
		panic(e)
	}

	orm.RunSyncdb("umbdynamicdb", false, true)
	fmt.Println("------------Setting schema Completed--------------------")
}
