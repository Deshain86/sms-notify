package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	Login         string = "matiopolskie"
	Password      string = "="
	ApiUrl        string = "https://webapi.allegro.pl/service.php"
	CountryCode   string = "1"
	LocalVersion  string = "1423185952"
	WebApiKey     string = ""
	UrlAllegro    string = `https://webapi.allegro.pl/service.php`
	SessionHandle string = "f8ef1b8ad951b38c5bd0b58152acfc8c9caadc09155a17//01_1"
)

func main() {
	log.Println("OK")
	doLogin()
	http.HandleFunc("/searchitem", searchItem)
	portStr := "3003"
	log.Printf("START PORT: %v\n", portStr)
	http.ListenAndServe(":"+portStr, nil)
}

type Request struct {
	Item string `json:"item"`
}

func searchItem(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	log.Printf("request\n%s\n", body)

	var dat Request
	if err := json.Unmarshal(body, &dat); err != nil {
		// panic(err)
		log.Println(err)
		return
	}

	log.Println(dat)
	x := getItem(dat.Item)
	val, err := json.Marshal(x)
	if err != nil {
		log.Println(err)
	}
	// log.Println(val)
	w.Write(val)
}

func doQueryAllSysStatus() {
	rq := ` <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="https://webapi.allegro.pl/service.php">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:DoQueryAllSysStatusRequest>
         <ser:countryId>` + CountryCode + `</ser:countryId>
         <ser:webapiKey>` + WebApiKey + `</ser:webapiKey>
      </ser:DoQueryAllSysStatusRequest>
   </soapenv:Body>
	</soapenv:Envelope>`
	res := SendHttpRequest(UrlAllegro, rq, "doQueryAllSysStatus")
	log.Println(res)

	b := bytes.NewBuffer([]byte(res))
	decoder := xml.NewDecoder(b)

	var inElement string
	var country string
	var verkey string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "item" {
				verkey = ""
				country = ""
			} else if inElement == "verKey" {
				var p string
				decoder.DecodeElement(&p, &se)
				verkey = p
			} else if inElement == "countryId" {
				var p string
				decoder.DecodeElement(&p, &se)
				country = p
			}
		case xml.EndElement:
			inElement = se.Name.Local
			if inElement == "item" {
				if country == "1" {
					LocalVersion = verkey
				}
			}
		}
	}
}
func getItem(item string) []Item {
	if SessionHandle == "" {
		return []Item{}
	}
	rq := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="https://webapi.allegro.pl/service.php">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:DoGetItemsListRequest>
         <ser:webapiKey>` + WebApiKey + `</ser:webapiKey>
         <ser:countryId>` + CountryCode + `</ser:countryId>
         <ser:filterOptions>
         <ser:item>
	      <ser:filterId>search</ser:filterId>
	      <ser:filterValueId>
	         <ser:item>` + item + `</ser:item>
	      </ser:filterValueId>
	   	</ser:item>
         </ser:filterOptions>
        <ser:resultSize>50</ser:resultSize>
        <ser:resultScope>3</ser:resultScope>
      </ser:DoGetItemsListRequest>
   </soapenv:Body>
</soapenv:Envelope>`
	// log.Println("\n\n" + rq + "\n\n")
	tmp := SendHttpRequest(UrlAllegro, rq, "DoGetItemsListRequest")
	b := bytes.NewBuffer([]byte(tmp))
	decoder := xml.NewDecoder(b)
	// log.Println(tmp)
	var inElement string

	var list []Item
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "item" {
				var p Item
				decoder.DecodeElement(&p, &se)
				list = append(list, p)
			}
		}
	}
	return list
}

type Item struct {
	ItemId     string             `xml:"itemId"`
	ItemTitle  string             `xml:"itemTitle"`
	PriceValue string             `xml:"priceValue"`
	PhotosInfo []PhotosInfoStruct `xml:"photosInfo>item"`
	EndingTime string             `xml:"endingTime"`
	TimeToEnd  string             `xml:"timeToEnd"`
}
type PhotosInfoStruct struct {
	PhotoUrl  string `xml:"photoUrl"`
	PhotoSize string `xml:"photoSize"`
}

func doLogin() {
	doQueryAllSysStatus()
	rq := ` <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="https://webapi.allegro.pl/service.php">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:DoLoginEncRequest>
         <ser:userLogin>` + Login + `</ser:userLogin>
         <ser:userHashPassword>` + Password + `</ser:userHashPassword>
         <ser:countryCode>` + CountryCode + `</ser:countryCode>
         <ser:webapiKey>` + WebApiKey + `</ser:webapiKey>
         <ser:localVersion>` + LocalVersion + `</ser:localVersion>
      </ser:DoLoginEncRequest>
   </soapenv:Body>
</soapenv:Envelope>`

	res := SendHttpRequest(UrlAllegro, rq, "doLoginEnc")
	b := bytes.NewBuffer([]byte(res))
	decoder := xml.NewDecoder(b)
	log.Println("RED ", res)
	var inElement string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "sessionHandlePart" {
				var p string
				decoder.DecodeElement(&p, &se)
				SessionHandle = p
			}
		}
	}
}

func SendHttpRequest(dst string, post string, method string) string {
	var headers = make([]string, 3)
	headers[0] = "Content-type:text/xml;charset=UTF-8"
	headers[1] = "SOAPaction:" + method

	client := &http.Client{}
	Method := ""
	Data := ""
	if post != "" {
		Method = "POST"
		Data = post
	} else {
		Method = "GET"
	}
	log.Println(Data)
	req, _ := http.NewRequest(Method, dst, bytes.NewReader([]byte(Data)))

	for _, x := range headers {
		if x != "" {
			y := strings.SplitN(x, ":", 2)
			req.Header.Add(y[0], y[1])
		}
	}
	response, _ := client.Do(req)
	defer response.Body.Close()
	resp, _ := ioutil.ReadAll(response.Body)

	return string(resp)
}
