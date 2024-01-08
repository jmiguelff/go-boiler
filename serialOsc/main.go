package main

import (
	"bufio"
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/hypebeast/go-osc/osc"
	"github.com/tarm/serial"
)

func main() {
	addrPtr := flag.String("addr", "localhost", "Destination IP OSC")
	portPtr := flag.Int("p", 7562, "Port OSC")
	flag.Parse()

	clientOsc := osc.NewClient(*addrPtr, *portPtr)

	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("reading on serial port: " + "/dev/ttyUSB0")

	r := bufio.NewReader(s)

	for {
		data, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		data = strings.TrimSuffix(data, "\n")
		log.Printf("[%x]", data)

		msgOsc := osc.NewMessage("/oscVal")
		// msgOsc.Append(data)
		// clientOsc.Send(msgOsc)

		dataSlice := strings.Split(data, " ")
		for _, v := range dataSlice {
			s, _ := strconv.ParseFloat(v, 32)
			msgOsc.Append(float32(s))
		}
		clientOsc.Send(msgOsc)
	}
}
