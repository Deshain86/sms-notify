package modem

import (
	"github.com/tarm/serial"
	"log"
	"strings"
	"time"
)

type Modem struct {
	comport  string
	bound    int
	instance *serial.Port
}

func New(comport string, bound int) *Modem {
	m := &Modem{comport: comport, bound: bound}
	c := &serial.Config{Name: comport, Baud: bound, ReadTimeout: time.Second}
	var err error
	m.instance, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

func (m *Modem) Send(number string, message string) {
	m.sendCommand("AT+CMGF=1\r", false)
	m.sendCommand("AT+CMGS=\""+number+"\"\r", false)
	x := m.sendCommand(message+string(26), true) // string 26 CTRL+Z
	log.Println("MESSAGE ", x)
}

func (m *Modem) sendCommand(message string, wait bool) string {
	m.instance.Flush()
	log.Println(message)
	_, err := m.instance.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 128)
	var loop int = 1
	if wait {
		loop = 10
	}
	var msg string
	var status string
	for i := 0; i < loop; i++ {
		n, _ := m.instance.Read(buf)
		if n > 0 {
			status = string(buf[:n])
			// log.Println(string(buf))
			msg += status
			if strings.HasSuffix(status, "OK\r\n") || strings.HasSuffix(status, "ERROR\r\n") {
				break
			}
		}
	}

	return msg
}
