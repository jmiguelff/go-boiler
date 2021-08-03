package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/tarm/serial"
)

func main() {
	fmt.Println("Starting simulator - open serial port")

	c := &serial.Config{Name: "/dev/ttymxc1", Baud: 38400}
	sfd, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 512)
	for {
		_, err := sfd.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buf)
	}
}

func sbs32WaitForHeader(s *serial.Port) int {
	// header we are expecting
	h := []byte{'\x00', '\xF1'}

	buf := make([]byte, 2)
	_, err := s.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	if bytes.Equal(buf, h) {
		s.Write([]byte("\xF2"))
		return 0
	} else {
		return -1
	}
}

func sbs32WaitForStartCmd(s *serial.Port) int {
	buf := make([]byte, 1)
	_, err := s.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	if buf[0] == '\xF2' {
		return 0
	} else {
		return -1
	}
}

//TODO: Implement cases for all possible commands
func sbs32WaitForCmd(s *serial.Port) int {
	buf := make([]byte, 1)
	_, err := s.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	if buf[0] == '\x49' {
		s.Write(buf)
	} else {
		return -1
	}

	_, err = s.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	if buf[0] == '\x42' {
		s.Write(buf)
	} else {
		return -1
	}

	return 0
}
