package models

import (
	"github.com/astaxie/beego"
	"encoding/xml"
	"github.com/astaxie/beego/orm"
	"fmt"
	"strconv"
	"strings"
	"github.com/astaxie/beego/httplib"
)

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

type Envelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Text    string   `xml:",chardata"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Body    struct {
		Text                 string `xml:",chardata"`
		GetINMainInfoRequest struct {
			Text    string `xml:",chardata"`
			Get     string `xml:"xmlns:get,attr"`
			TransId string `xml:"transId"`
			Msisdn  string `xml:"msisdn"`
		} `xml:"get:GetINMainInfoRequest"`
	} `xml:"soapenv:Body"`
}

type MainInfo struct {
	MaBalance     string `xml:"Body>GetINMainInfoResponse>maBalance"`
}

func init() {

}

var GetINMainInfoUrl = "http://10.147.114.5:8004/INServiceHandler/INBalance/GetINMainInfo_PS"

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

func GetUmb(msisdn string, mid string, sc string, cell string, regamtmn string, sms string) (v *Umb, err error) {

	beego.Info("-------------Model - umb.go------------------------")

	beego.Info("msisdn: ", msisdn)
	beego.Info("mid: ", mid)
	beego.Info("sc: ", sc)
	beego.Info("cell: ", cell)
	beego.Info("regamtmn: ", regamtmn)
	beego.Info("sms: ", sms)
	sms = strings.Replace(sms, " ", "%20", -1)

	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
    
	header := UmbHeader{MenuId: mid}
	
	
    err = o.Read(&header)
    if err == orm.ErrMultiRows {
        fmt.Printf("Returned Multi Rows Not One")
    }
    if err == orm.ErrNoRows {
        fmt.Printf("Not row found")
    }

    // Check if Header has BALANCE to be replaced
    if strings.Contains(header.MenuHeader, "%BALANCE%") {
    	mainInfo := MainInfo{}
    	balance, _ := GetINMainInfo(msisdn)
		xml.Unmarshal([]byte(balance), &mainInfo)
		header.MenuHeader = strings.Replace(header.MenuHeader, "%BALANCE%", mainInfo.MaBalance, -1)
    }

	Items :=[]Item{}
	var maps []orm.Params

    _, err = o.QueryTable("service_dyn_umb_menu").Filter("menu_id", mid).OrderBy("item_number").Values(&maps, "menu_id", "menu_detail_item", "item_number", "menu_next_id", "reg_amount", "unit", "formula", "keyword")
    if err == orm.ErrNoRows {
        fmt.Println("No records")
    } else if err == orm.ErrMissPK {
        fmt.Println("No Primary Key")
    } else {
		for i, v := range maps {
			detailItem, _ := v["MenuDetailItem"]
			itemNumber, _ := v["ItemNumber"].(string)
			nextId, _ := v["MenuNextId"].(string)
			index := strconv.Itoa(i)
			_ = index

			// If amount exist, change parameter regamtmn
			amount, _ := v["RegAmount"].(string)
			unit, _ := v["Unit"].(string)
			formula, _ := v["Formula"].(string)
			final_amount_str := "0"
			if amount != "" {
				if formula != "" {
					intAmount, _ := strconv.ParseInt(amount, 10, 0)
					intFormula , _ := strconv.ParseInt(formula, 10, 0)
					amount_str := (intAmount/intFormula)
					final_amount_str = strconv.FormatInt(amount_str, 10) + unit
					_ = final_amount_str
				} else {
					final_amount_str = amount
				}
			}

			// If keyword exist, change parameter sms
			keyword, _ := v["Keyword"].(string)
			if keyword != "" && strings.Contains(keyword, "%") {
				sms = keyword
			}

			// Insert Item 
			if str, ok := detailItem.(string); ok {
				if strings.Contains(str, "XXX") {
					str = strings.Replace(str, "XXX", final_amount_str, -1)
				}
				if mid == "ENDOFPSS" {
					sms = strings.Replace(sms, "__", "%20", -1)
					Items = append(Items, Item{
						str, itemNumber , "NONE", "0", "NEXT", "http://10.147.114.7:4080/PULLHandler/PullAPI_PS?origin=UMB&amp;sms=" + sms + "&amp;shortcode=123", "0",
					})
				} else {
					Items = append(Items, Item{
						str, itemNumber , "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=" + nextId + "&amp;regamtmn=" + final_amount_str + "&amp;bam=&amp;bbm=&amp;bcm=&amp;bdm=&amp;bem=&amp;sc=123&amp;sms=" + sms + "&amp;CELLID=" + cell + "&amp;param=", "0",
					})
				}
			} else {
				fmt.Println("q1q", ok)
			}
			
			
        }
    } 

    if regamtmn != "" {
    	header.MenuHeader = strings.Replace(header.MenuHeader, "XXX", regamtmn, -1)
    }
	v = &Umb{xml.Name{}, "", Event{"", "CPRespStatus", "SUCCESS", "0"}, Menu{"", "NONE", "0", header.MenuHeader, "", "superinternet", Items}}
	return v, nil
}

// Get data from backend API GetINMainInfo
func GetINMainInfo(id string) (str string, err2 error) {
	body := &Envelope{
		Soapenv:    "http://schemas.xmlsoap.org/soap/envelope/",
	}
	body.Body.GetINMainInfoRequest.Get = "http://www.example.org/GetINMainInfo/"
	body.Body.GetINMainInfoRequest.TransId = "111"
	body.Body.GetINMainInfoRequest.Msisdn = id

	req := httplib.Post(GetINMainInfoUrl)
	req.XMLBody(body)
	str, err := req.String()
	if err != nil {
		beego.Info(err)
	}
	return str, nil
}