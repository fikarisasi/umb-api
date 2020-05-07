package models

import "encoding/xml"

type Umb struct {
	XMLName xml.Name `xml:"umb"`
	Text    string   `xml:",chardata"`
	Event   Event    `xml:"event"`
	Menu    Menu     `xml:"menu"`
}

type Event struct {
	Text       string `xml:",chardata"`
	Opcode     string `xml:"opcode,attr"`
	Status     string `xml:"status,attr"`
	Statuscode string `xml:"statuscode,attr"`
}

type Menu struct {
	Text        string `xml:",chardata"`
	Tarifftype  string `xml:"tarifftype,attr"`
	Tariffrate  string `xml:"tariffrate,attr"`
	Menuheader1 string `xml:"menuheader1"`
	Menuheader2 string `xml:"menuheader2"`
	Menuname    string `xml:"menuname"`
	Item        []Item `xml:"item"`
}

type Item struct {
	Text         string `xml:",chardata"`
	Number       string `xml:"number,attr"`
	Tarifftype   string `xml:"tarifftype,attr"`
	Tariffrate   string `xml:"tariffrate,attr"`
	Action       string `xml:"action,attr"`
	URL          string `xml:"url,attr"`
	DeliveryMode string `xml:"delivery_mode,attr"`
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
	Items := []Item{
		Item{
			"OK", "1", "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=POSTPAID_SS1&amp;regamtmn=0&amp;bam=&amp;bbm=&amp;bcm=&amp;bdm=&amp;bem=&amp;sc=123&amp;sms=&amp;CELLID=999999&amp;param=", "0",
		},
		Item{
			"Benefit Postpaid", "2", "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=POSTPAID_SS2&amp;regamtmn=0&amp;bam=&amp;bbm=&amp;bcm=&amp;bdm=&amp;bem=&amp;sc=123&amp;sms=&amp;CELLID=999999&amp;param=", "0",
		},
		Item{
			"Status", "3", "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=POSTPAID_SS3&amp;regamtmn=0&amp;bam=&amp;bbm=&amp;bcm=&amp;bdm=&amp;bem=&amp;sc=123&amp;sms=&amp;CELLID=999999&amp;param=", "0",
		},
		Item{
			"Info", "4", "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=INFO_SS1&amp;regamtmn=0&amp;bam=&amp;bbm=&amp;bcm=&amp;bdm=&amp;bem=&amp;sc=123&amp;sms=&amp;CELLID=999999&amp;param=", "0",
		},
	}
	v = &Umb{xml.Name{}, "", Event{"", "CPRespStatus", "SUCCESS", "0"}, Menu{"", "NONE", "0", "Berlangganan Postpaid, Lebih Banyak Untung !", "", "superinternet", Items}}
	return v, nil
}
