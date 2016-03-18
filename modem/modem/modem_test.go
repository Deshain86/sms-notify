package modem_test

import (
	modem "../modem"
	"log"
	"testing"
)

func Test_modem(t *testing.T) {
	log.Println("OK")
	m := modem.New("/dev/ttyUSB0", 115200)
	// m.ReadAll()
	msg := "Gdzie jest foodtruck?"
	m.Send("+48691157964", msg)
	// m.Send("+48691563056", msg)
	// m.Send("+48791207146", msg)
}
