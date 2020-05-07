package models

import (
	"encoding/xml"
)

type Umb struct {
	XMLName xml.Name `xml:"umb"`
	Text    string   `xml:",chardata"`
	Event   struct {
		Text       string `xml:",chardata"`
		Opcode     string `xml:"opcode,attr"`
		Status     string `xml:"status,attr"`
		Statuscode string `xml:"statuscode,attr"`
	} `xml:"event"`
	Menu struct {
		Text        string `xml:",chardata"`
		Tarifftype  string `xml:"tarifftype,attr"`
		Tariffrate  string `xml:"tariffrate,attr"`
		Menuheader1 string `xml:"menuheader1"`
		Menuheader2 string `xml:"menuheader2"`
		Menuname    string `xml:"menuname"`
		Item        []struct {
			Text         string `xml:",chardata"`
			Number       string `xml:"number,attr"`
			Tarifftype   string `xml:"tarifftype,attr"`
			Tariffrate   string `xml:"tariffrate,attr"`
			Action       string `xml:"action,attr"`
			URL          string `xml:"url,attr"`
			DeliveryMode string `xml:"delivery_mode,attr"`
		} `xml:"item"`
	} `xml:"menu"`
}

func init() {

}

// GetUmbById retrieves Umb by Id. Returns error if
// Id doesn't exist
func GetUmbById(id int64) (v *Umb, err error) {
	// o := orm.NewOrm()
	// v = &Umb{Id: id}
	// if err = o.QueryTable(new(Umb)).Filter("Id", id).RelatedSel().One(v); err == nil {
	// 	return v, nil
	// }
	// return nil, err
	v = &Umb{}
	return v, nil
}
