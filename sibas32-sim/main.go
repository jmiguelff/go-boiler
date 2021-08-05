package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/tarm/serial"
)

type state int

const (
	WAITFORHEADER state = iota
	WAITFORCMD
	REQCONNECTORBYID
	REQBAUDRATECHANGE
	CHANGEBAUD
	REQSWVERSION
	ERROR
)

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

func main() {
	fmt.Println("Starting simulator - open serial port")

	c := &serial.Config{Name: "/dev/ttymxc0", Baud: 38400}
	sfd, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalln(err)
	}

	var headerOk int = 0
	var selCmd cmd = NA
	var cmdOk int = 0
	var serialOk int = 0

	s := WAITFORHEADER
	for {
		switch s {
		case WAITFORHEADER:
			if headerOk == 1 {
				s = WAITFORCMD
			} else {
				s = WAITFORHEADER
			}
		case WAITFORCMD:
			if selCmd == SETBAUDRATE {
				s = REQBAUDRATECHANGE
			} else if selCmd == GETCONNECTORBYID {
				s = REQCONNECTORBYID
			} else if selCmd == GETSWVERSION {
				s = REQSWVERSION
			} else if selCmd == NA && headerOk == 1 {
				s = WAITFORCMD
			} else if selCmd == NA && headerOk == 0 {
				s = WAITFORHEADER
			} else {
				s = ERROR
			}
		case REQCONNECTORBYID:
			if cmdOk == 0 {
				s = REQCONNECTORBYID
			} else {
				s = WAITFORHEADER
			}
		case REQBAUDRATECHANGE:
			if cmdOk == 0 {
				s = REQBAUDRATECHANGE
			} else {
				s = CHANGEBAUD
			}
		case REQSWVERSION:
			if cmdOk == 0 {
				s = REQSWVERSION
			} else {
				s = WAITFORHEADER
			}
		case CHANGEBAUD:
			if serialOk == 1 {
				s = WAITFORHEADER
			} else {
				log.Fatalln("Fail to change baudrate")
			}
		default:
			s = ERROR
		}

		switch s {
		case WAITFORHEADER:
			cmdOk = 0
			selCmd = NA
			err = sbs32WaitForHeader(sfd)
			if err != nil {
				headerOk = 0
				log.Println(err)
			} else {
				headerOk = 1
			}
		case WAITFORCMD:
			selCmd, err = sbs32WaitForCmd(sfd)
			if err != nil {
				headerOk = 0
				selCmd = NA
				log.Println(err)
			}
		case REQCONNECTORBYID:
			adr, err := sbs32WaitForFPId(sfd)
			if err != nil {
				headerOk = 0
				selCmd = NA
				log.Println(err)
			} else {
				log.Println(adr)
				err = sbs32WaitForFooter(sfd)
				if err != nil {
					headerOk = 0
					selCmd = NA
					log.Println(err)
				} else {
					err = sbs32SendConnectorVal(sfd)
					if err != nil {
						headerOk = 0
						selCmd = NA
						log.Println(err)
					} else {
						cmdOk = 1
						headerOk = 0
					}
				}
			}
		case REQBAUDRATECHANGE:
			err = sbs32EchoPayloadBaud(sfd)
			if err != nil {
				headerOk = 0
				selCmd = NA
				log.Println(err)
			} else {
				err = sbs32WaitForFooterBaud(sfd)
				if err != nil {
					headerOk = 0
					selCmd = NA
					log.Println(err)
				} else {
					cmdOk = 1
					headerOk = 0
				}
			}
		case REQSWVERSION:
			err = sbs32SendVersion(sfd)
			if err != nil {
				headerOk = 0
				selCmd = NA
				log.Println(err)
			} else {
				cmdOk = 1
				headerOk = 0
			}
		case CHANGEBAUD:
			serialOk = 0
			sfd.Close()

			c.Baud = 115200
			sfd, err = serial.OpenPort(c)
			if err != nil {
				log.Fatalln(err)
			} else {
				serialOk = 1
			}
		case ERROR:
			log.Fatalln("Fatal error - exit")
		default:
			log.Fatalln("Failed to default state on actions switch/case")
		}

	}
}

func sbs32WaitForHeader(s *serial.Port) error {
	// header we are expecting
	h := []byte{'\x00', '\xF1'}

	buf := make([]byte, 2)
	_, err := s.Read(buf)
	if err != nil {
		return err
	}

	if bytes.Equal(buf, h) {
		s.Write([]byte("\xF2"))
		return nil
	} else {
		return errors.New("incorrect header received")
	}
}

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
	buf := make([]byte, 2)
	_, err := s.Read(buf)
	if err != nil {
		return NA, err
	}

	var aux int

	switch buf[0] {
	case '\x48':
		s.Write(buf[:1])
		return SETBAUDRATE, nil // Set baudrate
	case '\x49':
		s.Write(buf[:1])
		aux = 1 // Read connector value by FP id or get FP name (req next byte)
	case '\x4a':
		s.Write(buf[:1])
		aux = 2 // Get connector value by FP name or Get connector address (req next byte)
	case '\x4b':
		s.Write(buf[:1])
		aux = 3 // Get connector id (req next byte)
	case '\x6a':
		s.Write(buf[:1])
		return GETNBROFP, nil // Get number FP
	case '\x72':
		s.Write(buf[:1])
		return GETFPACCESSRIGHTS, nil // Get FP access rights
	case '\x73':
		s.Write(buf[:1])
		return GETSWVERSION, nil // Get version
	default:
		return NA, errors.New("invalid command")
	}

	// Get second byte
	n, err := s.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	if n > 1 {
		return NA, errors.New("received more data than expected")
	}

	switch buf[0] {
	case '\x41':
		if aux == 1 {
			s.Write(buf[:1])
			return GETFPNAME, nil // Get FP name
		} else {
			return NA, errors.New("invalid command (2nd byte)")
		}
	case '\x42':
		if aux == 1 {
			s.Write(buf[:1])
			return GETCONNECTORBYID, nil // Read connector value by id
		} else if aux == 2 {
			s.Write(buf[:1])
			return GETCONNECTORBYNAME, nil // Get connector value by FP name
		} else if aux == 3 {
			s.Write(buf[:1])
			return GETCONNECTORID, nil // Get connector id
		} else {
			return NA, errors.New("invalid command (2nd byte)")
		}
	case '\x44':
		if aux == 2 {
			s.Write(buf[:1])
			return GETCONNECTORADDRESS, nil // Get connector address
		} else {
			return NA, errors.New("invalid command (2nd byte)")
		}
	default:
		return NA, errors.New("invalid command (2nd byte)")
	}
}

func sbs32WaitForFPId(s *serial.Port) ([]byte, error) {
	id := make([]byte, 5)

	for i := 0; i < len(id); i++ {
		buf := make([]byte, 1)
		_, err := s.Read(buf)
		if err != nil {
			return id, err
		} else {
			s.Write(buf)
			id = append(id, buf[0])
		}
	}
	return id, nil
}

func sbs32WaitForFooter(s *serial.Port) error {
	buf := make([]byte, 1)
	_, err := s.Read(buf)
	if err != nil {
		return err
	}

	if buf[0] == '\x4f' {
		return nil
	} else {
		return errors.New("invalid footer received")
	}
}

func sbs32SendConnectorVal(s *serial.Port) error {
	data := []byte{
		'\x45', '\x4E', '\x41', '\x42', '\x44', '\x49', '\x41', '\x33', '\x2E', '\x4F', '\x55', '\x54', '\x33', '\x36', '\x30', '\x30',
		'\x30', '\x30', '\x30', '\x30', '\x4F', '\x39', '\x45', '\x30', '\x30', '\x45', '\x42', '\x32', '\x43', '\x30', '\x30', '\x32',
		'\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x38', '\x46', '\x33', '\x20', '\x20', '\x20', '\x20', '\x20', '\x20', '\x20',
		'\x20', '\x20', '\x20', '\x32', '\x36', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x31', '\x30', '\x30', '\x30', '\x30',
		'\x30', '\x30', '\x30'}

	_, err := s.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func sbs32SendVersion(s *serial.Port) error {
	data := []byte{
		'\x4F', '\x32', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30',
		'\x30', '\x43', '\x43', '\x55', '\x20', '\x66', '\x6F', '\x72', '\x20', '\x43', '\x50', '\x41', '\x20', '\x53', '\x82', '\x72',
		'\x69', '\x65', '\x20', '\x34', '\x30', '\x30', '\x30', '\x0A', '\x50', '\x47', '\x4D', '\x20', '\x0A', '\x50', '\x47', '\x4D',
		'\x20', '\x76', '\x65', '\x72', '\x73', '\x69', '\x6F', '\x6E', '\x3A', '\x20', '\x32', '\x2E', '\x30'}

	_, err := s.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func sbs32EchoPayloadBaud(s *serial.Port) error {
	buf := make([]byte, 1)
	for i := 0; i < 4; i++ {
		_, err := s.Read(buf)
		if err != nil {
			return err
		}
		s.Write(buf)
	}
	return nil
}

func sbs32WaitForFooterBaud(s *serial.Port) error {
	footer := []byte{'\x4F', '\xFC', '\xFC', '\xFC', '\xFC', '\xFC'}
	buf := make([]byte, 6)
	_, err := s.Read(buf)
	if err != nil {
		return err
	}

	if bytes.Equal(buf, footer) {
		return nil
	} else {
		return errors.New("footer does not match")
	}
}
