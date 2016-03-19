package main

import (
	"./modem"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

var m *modem.Modem
var db *sql.DB

func main() {
	log.Println("OK")
	var err error
	db, err = sql.Open("mysql", "root:root@/sms")
	if err != nil {
		return
	}

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
		phone = "48" + dat.Phone
	} else if len(dat.Phone) == 11 {
		phone = dat.Phone
	} else {
		return
	}
	id := addToDB(0, phone, dat.Message)
	err := m.Send("+"+phone, dat.Message)
	if err != nil {
		log.Println(err)
	} else {
		updateToDb(id)
	}
	log.Printf("%#v", dat)
}
