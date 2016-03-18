package main

import (
	"./modem"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var m *modem.Modem

func main() {
	log.Println("OK")

	m = modem.New("/dev/ttyUSB0", 115200)

	http.HandleFunc("/sendsms", sendSms)
	portStr := "3002"
	log.Printf("START PORT: %v\n", portStr)
	http.ListenAndServe(":"+portStr, nil)
}

type Request struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func sendSms(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	log.Printf("request\n%s\n", body)

	var dat Request
	if err := json.Unmarshal(body, &dat); err != nil {
		// panic(err)
		log.Println(err)
		return
	}

	if len(dat.Message) == 0 {
		return
	}
	var phone string
	if len(dat.Phone) == 9 {
		phone = "+48" + dat.Phone
	} else if len(dat.Phone) == 12 {
		phone = dat.Phone
	} else {
		return
	}
	m.Send(phone, dat.Message)
	log.Printf("%#v", dat)
}
