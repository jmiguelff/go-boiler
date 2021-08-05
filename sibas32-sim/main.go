package main

import (
	"bytes"
	"errors"
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

type cmd int

const (
	SETBAUDRATE cmd = iota
	GETNBROFP
	GETFPACCESSRIGHTS
	GETSWVERSION
	GETFPNAME
	GETCONNECTORBYID
	GETCONNECTORBYNAME
	GETCONNECTORID
	GETCONNECTORADDRESS
	NA
)

// cmd list
// 0: Set baudrate
// 1: Get number of FP
// 2: Get FP access rights
// 3: Get sw version
// 4: Get FP name
// 5: Read connector value by id
// 6: Read connector value by FP name
// 7: Get connector id
// 8: Get connector address
func sbs32WaitForCmd(s *serial.Port) (cmd, error) {
	// Get first byte
	buf := make([]byte, 1)
	_, err := s.Read(buf)
	if err != nil {
		return NA, err
	}

	var aux int

	switch buf[0] {
	case '\x48':
		s.Write(buf)
		return SETBAUDRATE, nil // Set baudrate
	case '\x49':
		s.Write(buf)
		aux = 1 // Read connector value by FP id or get FP name (req next byte)
	case '\x4a':
		s.Write(buf)
		aux = 2 // Get connector value by FP name or Get connector address (req next byte)
	case '\x4b':
		s.Write(buf)
		aux = 3 // Get connector id (req next byte)
	case '\x6a':
		s.Write(buf)
		return GETNBROFP, nil // Get number FP
	case '\x72':
		s.Write(buf)
		return GETFPACCESSRIGHTS, nil // Get FP access rights
	case '\x73':
		s.Write(buf)
		return GETSWVERSION, nil // Get version
	default:
		return NA, errors.New("invalid command")
	}

	// Get second byte
	_, err = s.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	switch buf[0] {
	case '\x41':
		if aux == 1 {
			s.Write(buf)
			return GETFPNAME, nil // Get FP name
		} else {
			return NA, errors.New("invalid command (2nd byte)")
		}
	case '\x42':
		if aux == 1 {
			s.Write(buf)
			return GETCONNECTORBYID, nil // Read connector value by id
		} else if aux == 2 {
			s.Write(buf)
			return GETCONNECTORBYNAME, nil // Get connector value by FP name
		} else if aux == 3 {
			s.Write(buf)
			return GETCONNECTORID, nil // Get connector id
		} else {
			return NA, errors.New("invalid command (2nd byte)")
		}
	case '\x44':
		if aux == 2 {
			s.Write(buf)
			return GETCONNECTORADDRESS, nil // Get connector address
		} else {
			return NA, errors.New("invalid command (2nd byte)")
		}
	default:
		return NA, errors.New("invalid command (2nd byte)")
	}
}

// Change it to look for the available commands
