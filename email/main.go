package main

import (
	"log"
	"strings"
	"time"

	"gopkg.in/mvader/go-imapreader.v1"
)

const (
	Addr = "maxhumor.nazwa.pl"
	User = "angrygophers@firus.pl"
	Pass = "Angrygophers007"
	MBox = "INBOX"
)

type dataset struct {
	Addr string
	User string
	Pass string
}

var activeReaders map[string]imapreader.Reader

func main() {

	var datalist []dataset
	datalist = append(datalist, dataset{Addr: Addr, User: User, Pass: Pass})
	
	activeReaders := make(map[string]imapreader.Reader)

	for _, set := range datalist {

		if _, ok := activeReaders[set.Addr+set.User]; !ok {

			r, err := imapreader.NewReader(imapreader.Options{
				Addr:     Addr,
				Username: User,
				Password: Pass,
				TLS:      true,
				Timeout:  60 * time.Second,
				MarkSeen: true,
			})
			if err != nil {
				log.Print(err)
				continue
			}
			activeReaders[set.Addr+set.User] = r

			if err := r.Login(); err != nil {
				panic(err)
			}
			defer r.Logout()

			// Search for all the emails in "all mail" that are unseen
			// read the docs for more search filters

			imapFolder := ""
			if strings.Contains(Addr, "gmail.com") {
				imapFolder = imapreader.GMailAllMail
			} else {
				imapFolder = imapreader.GMailInbox
			}

			messages, err := r.List(imapFolder, imapreader.SearchAll) //imapreader.SearchUnseen)
			if err != nil {
				panic(err)
			}

			for _, x := range messages {
				receiveDate, err := x.Header.Date()
				if err != nil {
					panic(err)
				}
				log.Printf("%#v", receiveDate)
				log.Printf("%#v", x.Header.Get("From"))
				log.Printf("%#v", x.Header.Get("Subject"))
			}

		}

	}

	// do stuff with messages
}
