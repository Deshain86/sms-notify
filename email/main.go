package main

import (
	"log"
	"strings"
	"sync"
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

var datalist []dataset

type readersStruct struct {
	Locked        *sync.Mutex
	ActiveReaders map[string]imapreader.Reader
}

type Ticker func(time.Time)

var readers readersStruct

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	readers.ActiveReaders = make(map[string]imapreader.Reader)
	readers.Locked = &sync.Mutex{}
	go setTicker(manageReceivers)
	<- time.After(3 * time.Second)
	go setTicker(readEmails)
	manageReceivers(time.Now())
	<- time.After(3 * time.Second)
	readEmails(time.Now())
	for {
	}
}

func readEmails(now time.Time) {
	readers.Locked.Lock()
	for _, r := range readers.ActiveReaders {
		// Search for all the emails in "all mail" that are unseen
		// read the docs for more search filters

		imapFolder := ""
		if strings.Contains(Addr, "gmail.com") {
			imapFolder = imapreader.GMailAllMail
		} else {
			imapFolder = imapreader.GMailInbox
		}
		log.Print(imapFolder)
		messages, err := r.List(imapFolder, imapreader.SearchAll) //imapreader.SearchUnseen)
		if err != nil {
			log.Print(err)
			continue
		}

		for _, x := range messages {
			receiveDate, err := x.Header.Date()
			if err != nil {
				log.Print(err)
				continue
			}
			log.Printf("%#v", receiveDate)
			log.Printf("%#v", x.Header.Get("From"))
			log.Printf("%#v", x.Header.Get("Subject"))
		}
	}
	readers.Locked.Unlock()
}

func setTicker(what Ticker) {
	c := time.Tick(15 * time.Second)
	for now := range c {
		what(now)
	}
}

func manageReceivers(now time.Time) {
	readers.ActiveReaders = make(map[string]imapreader.Reader)
	datalist = datalist[:0]
	//TUTAJ DB
	datalist = append(datalist, dataset{Addr: Addr, User: User, Pass: Pass})
	readers.Locked.Lock()
	for _, set := range datalist {
		log.Printf("Time: %s Set: %#v", now, set)
		if _, ok := readers.ActiveReaders[set.Addr+set.User]; !ok {

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

			if err := r.Login(); err != nil {
				log.Print(err)
				continue
			}
			readers.ActiveReaders[set.Addr+set.User] = r
//			defer r.Logout()
		}
	}
	readers.Locked.Unlock()

}
