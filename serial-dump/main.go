package main

import (
	"log"

	"github.com/tarm/serial"
)

func main() {
	// Open serial port with 115200 baud, 1 stop bit, no parity (8N1)
	c := &serial.Config{
		Name:     "/dev/ttymxc0",
		Baud:     115200,
		Size:     8,                 // 8 data bits
		StopBits: serial.Stop1,      // 1 stop bit
		Parity:   serial.ParityNone, // no parity
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// Create 2000 bytes of data (all 'A's here)
	data := make([]byte, 2000)
	for i := range data {
		data[i] = 'A'
	}

	// Single system call to write all 2000 bytes
	n, err := s.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sent %d bytes\n", n)
}
