package models

import (
	"github.com/astaxie/beego"
	"encoding/xml"
	"github.com/astaxie/beego/orm"
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"github.com/astaxie/beego/httplib"
)

type Umb struct {
	XMLName xml.Name `xml:"umb"`
	Text    string   `xml:",chardata"`
	Event   Event    `xml:"event"`
	Menu    *Menu     `xml:",omitempty"`
	Result  *Result   `xml:",omitempty"`
}

type Event struct {
	Text       string `xml:",chardata"`
	Opcode     string `xml:"opcode,attr"`
	Status     string `xml:"status,attr"`
	Statuscode string `xml:"statuscode,attr"`
}

type Menu struct {
	XMLName 		xml.Name `xml:"menu"`
	Tarifftype  	string `xml:"tarifftype,attr,omitempty"`
	Tariffrate  	string `xml:"tariffrate,attr,omitempty"`
	Menuheader1 	string `xml:"menuheader1,omitempty"`
	Menuheader2 	string `xml:"menuheader2"`
	Menuname    	string `xml:"menuname,omitempty"`
	Item        	[]Item `xml:"item"`
}

type Result struct {
	XMLName 		xml.Name `xml:"result"`
	Resultdata    	string `xml:"resultdata,omitempty"`
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

type errorRes struct {
	Category   	string `json:"category"`
    Message 	error `json:"message"`
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

type CRSInfo struct {
	Desc string `xml:"DESC"` 
    Name [] string `xml:"ATTRIBUTES>KEY>NAME"` 
	Value [] string `xml:"ATTRIBUTES>KEY>VALUE"`
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
	// Items := []Item{
	// 	Item{
	// 		"OK", "1", "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=POSTPAID_SS1&regamtmn=0&bam=&bbm=&bcm=&bdm=&bem=&sc=123&sms=&CELLID=999999&param=", "0",
	// 	},
	// 	Item{
	// 		"Benefit Postpaid", "2", "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=POSTPAID_SS2&regamtmn=0&bam=&bbm=&bcm=&bdm=&bem=&sc=123&sms=&CELLID=999999&param=", "0",
	// 	},
	// 	Item{
	// 		"Status", "3", "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=POSTPAID_SS3&regamtmn=0&bam=&bbm=&bcm=&bdm=&bem=&sc=123&sms=&CELLID=999999&param=", "0",
	// 	},
	// 	Item{
	// 		"Info", "4", "NONE", "0", "NEXT", "http://10.147.114.7:4080/UMB/Menu?mid=INFO_SS1&regamtmn=0&bam=&bbm=&bcm=&bdm=&bem=&sc=123&sms=&CELLID=999999&param=", "0",
	// 	},
	// }
	// v = &Umb{xml.Name{}, "", Event{"", "CPRespStatus", "SUCCESS", "0"}, Menu{"NONE", "0", "Berlangganan Postpaid, Lebih Banyak Untung !", "", "superinternet", Items}}
	// return v, nil
	v = &Umb{
		XMLName: xml.Name{},
		Event: Event{"", "CPRespStatus", "SUCCESS", "0"},
	}
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
	xmlResult := true
    o := orm.NewOrm()
    o.Using("default")

    urlhost := "http://umbmenu-oc.office.corp.indosat.com/"

    // Checking Postpaid SS if any wrong format
    if strings.Contains(mid, "_FALSE") {
    	xmlResult = false
    }
    
	header := UmbHeader{MenuId: mid}
	msisdnNIK := ""
	msisdnNOKK := ""
	msisdnRegistered := ""
	
    err = o.Read(&header)
    if err == orm.ErrMultiRows {
        fmt.Printf("Returned Multi Rows Not One")
    }
    if err == orm.ErrNoRows {
        fmt.Printf("Not row found")
    }
    if mid == "" {
    	header.MenuHeader = "Infomation not found"
    }

    // Check if Header has BALANCE to be replaced
    if strings.Contains(header.MenuHeader, "%BALANCE%") {
    	mainInfo := MainInfo{}
    	balance, _ := GetINMainInfo(msisdn)
		xml.Unmarshal([]byte(balance), &mainInfo)
		header.MenuHeader = strings.Replace(header.MenuHeader, "%BALANCE%", mainInfo.MaBalance, -1)
    }

    // Check if Header has NIK & NOKK to be replaced
    if strings.Contains(header.MenuHeader, "%NIK%") && strings.Contains(header.MenuHeader, "%NOKK%") {
    	crsInfo := CRSInfo{Name: [] string{} , Value: [] string{}}
    	crsValue, _ := CRSHandler(msisdn)
	    xml.Unmarshal([]byte(crsValue), &crsInfo)

	    // Check if MSISDN registered or not
	    if crsInfo.Desc == "REGISTERED" {
	    	msisdnRegistered = crsInfo.Desc
	    }
		
		for i := range crsInfo.Name {
		  	if(crsInfo.Name[i] == "NIK"){
		  		msisdnNIK = crsInfo.Value[i]
		  		header.MenuHeader = strings.Replace(header.MenuHeader, "%NIK%", msisdnNIK, -1)
		  	}
	        if(crsInfo.Name[i] == "NOKK"){
	        	msisdnNOKK = crsInfo.Value[i]
	        	header.MenuHeader = strings.Replace(header.MenuHeader, "%NOKK%", msisdnNOKK, -1)
	        }
		}
    }

    // Check if Header has | 
    if strings.Contains(header.MenuHeader, "|") {
    	subPipeSring := strings.Split(header.MenuHeader, "|")
    	header.MenuHeader = subPipeSring[0]
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

			final_sms := sms
			// If keyword exist, change parameter sms
			keyword, _ := v["Keyword"].(string)
			if keyword != "" && strings.Contains(keyword, "%") {
				final_sms = keyword
			}

			// Insert Item 
			if str, ok := detailItem.(string); ok {
				if strings.Contains(str, "XXX") {
					str = strings.Replace(str, "XXX", final_amount_str, -1)
				}
				if msisdnRegistered != "" && strings.Contains(final_sms, "%20DESC%20") {
					final_sms = strings.Replace(final_sms, "DESC", msisdnRegistered, -1)
				}
				if msisdnNIK != "" && strings.Contains(final_sms, "%20NIK%20") {
					final_sms = strings.Replace(final_sms, "NIK", msisdnNIK, -1)
				}
				if msisdnNOKK != "" && strings.Contains(final_sms, "%20NOKK%20") {
					final_sms = strings.Replace(final_sms, "NOKK", msisdnNOKK, -1)
				}
				if mid == "ENDOFPSS" {
					sms = strings.Replace(sms, "__", "%20", -1)
					Items = append(Items, Item{
						str, itemNumber , "NONE", "0", "NEXT", urlhost + "PULLHandler/PullAPI_PS?origin=UMB&sms=" + final_sms + "&shortcode=123", "0",
					})
				} else {
					Items = append(Items, Item{
						str, itemNumber , "NONE", "0", "NEXT", urlhost + "UMB/Menu?mid=" + nextId + "&regamtmn=" + final_amount_str + "&bam=&bbm=&bcm=&bdm=&bem=&sc=123&sms=" + final_sms + "&CELLID=" + cell + "&param=", "0",
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
    if len(Items) == 0 {
    	xmlResult = false
    }
	// v = &Umb{xml.Name{}, "", Event{"", "CPRespStatus", "SUCCESS", "0"}, Menu{"", "NONE", "0", header.MenuHeader, "", "superinternet", Items}}
	v = &Umb{
		XMLName: xml.Name{},
		Event: Event{"", "CPRespStatus", "SUCCESS", "0"},
	}
	if (xmlResult) {
		v.Menu = &Menu{
			// XMLName: xml.Name{ Local: "menu" },
			Tarifftype: "NONE",
			Tariffrate: "0",
			Menuheader1: header.MenuHeader,
			Menuheader2: "",
			Menuname: "superinternet",
			Item: Items,
		}
	} else {
		v.Result = &Result{
			// XMLName: xml.Name{ Local: "result" },
			Resultdata: header.MenuHeader,
		}
	}
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
		ginmiErrRes := &errorRes {
			Category: "GetINMainInfo API Error",
			Message: err,
		}
		ginmiErrResJSON, _ := json.Marshal(ginmiErrRes)
    	beego.Error(string(ginmiErrResJSON))
		beego.Error(ginmiErrRes)
	}
	return str, nil
}

// Get data from backend API CRSHandler
func CRSHandler(id string) (str string, err2 error) {
	body := &Envelope{
		Soapenv:    "http://schemas.xmlsoap.org/soap/envelope/",
	}
	body.Body.GetINMainInfoRequest.Get = "http://www.example.org/GetINMainInfo/"
	body.Body.GetINMainInfoRequest.TransId = "111"
	body.Body.GetINMainInfoRequest.Msisdn = id

	var CRSHandlerUrl = "http://10.34.36.68:8080/SelfcareRegistrationStatus/" + id + "?clientId=USSD&cred=d8dd050b3d326872ef50301047b04125"

	req := httplib.Get(CRSHandlerUrl)
	req.XMLBody(body)
	str, err := req.String()
	if err != nil {
		crsErrRes := &errorRes {
			Category: "SelfcareRegistrationStatus API Error",
			Message: err,
		}
		crsErrResJSON, _ := json.Marshal(crsErrRes)
    	beego.Error(string(crsErrResJSON))
		beego.Error(crsErrRes)
	}
	return str, nil
}