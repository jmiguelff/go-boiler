package main

import (
	"bufio"
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

	c := &serial.Config{Name: "/dev/ttymxc2", Baud: 38400}
	sfd, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalln(err)
	}
	r := bufio.NewReader(sfd)

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
			err = sbs32WaitForHeader(sfd, r)
			if err != nil {
				headerOk = 0
				log.Println(err)
			} else {
				headerOk = 1
			}
		case WAITFORCMD:
			selCmd, err = sbs32WaitForCmd(sfd, r)
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
			err = sbs32VersionQuery(sfd, r)
			if err != nil {
				log.Println(err)
				cmdOk = -1
				headerOk = 0
				selCmd = NA
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

func sbs32WaitForHeader(s *serial.Port, r *bufio.Reader) error {
	// header we are expecting
	h := []byte{'\x00', '\xF1'}

	// wait for first header byte
	for {
		buf, err := r.ReadByte()
		if err != nil {
			return err
		}
		if buf == h[0] {
			log.Printf("Received first header byte [%x]\n", buf)
			break
		} else {
			continue
		}
	}

	// get second header byte
	buf, err := r.ReadByte()
	if err != nil {
		return err
	}

	if buf == h[1] {
		log.Printf("Received second header byte [%x]\n", buf)
	} else {
		return errors.New("second header byte does not match prrotocol")
	}

	// send ack byte
	_, err = s.Write([]byte{'\xF2'})
	if err != nil {
		return err
	}

	// get ack echo byte
	buf, err = r.ReadByte()
	if err != nil {
		return err
	}

	if buf != '\xF2' {
		return errors.New("ack echo byte does not match")
	}

	return nil
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
func sbs32WaitForCmd(s *serial.Port, r *bufio.Reader) (cmd, error) {
	// Get first byte
	buf, err := r.ReadByte()
	if err != nil {
		return NA, err
	}

	// debug print
	log.Printf("Received first byte command [%x]\n", buf)

	var aux int
	switch buf {
	case '\x48':
		s.Write([]byte{buf})
		return SETBAUDRATE, nil // Set baudrate
	case '\x49':
		s.Write([]byte{buf})
		aux = 1 // Read connector value by FP id or get FP name (req next byte)
	case '\x4a':
		s.Write([]byte{buf})
		aux = 2 // Get connector value by FP name or Get connector address (req next byte)
	case '\x4b':
		s.Write([]byte{buf})
		aux = 3 // Get connector id (req next byte)
	case '\x6a':
		s.Write([]byte{buf})
		return GETNBROFP, nil // Get number FP
	case '\x72':
		s.Write([]byte{buf})
		return GETFPACCESSRIGHTS, nil // Get FP access rights
	case '\x73':
		s.Write([]byte{buf})
		return GETSWVERSION, nil // Get version
	default:
		return NA, errors.New("invalid command")
	}

	// get second byte
	buf, err = r.ReadByte()
	if err != nil {
		return NA, err
	}

	// debug
	log.Printf("Received second byte command [%x]\n", buf)

	switch buf {
	case '\x41':
		if aux == 1 {
			s.Write([]byte{buf})
			return GETFPNAME, nil // Get FP name
		} else {
			return NA, errors.New("invalid command (2nd byte)")
		}
	case '\x42':
		if aux == 1 {
			s.Write([]byte{buf})
			return GETCONNECTORBYID, nil // Read connector value by id
		} else if aux == 2 {
			s.Write([]byte{buf})
			return GETCONNECTORBYNAME, nil // Get connector value by FP name
		} else if aux == 3 {
			s.Write([]byte{buf})
			return GETCONNECTORID, nil // Get connector id
		} else {
			return NA, errors.New("invalid command (2nd byte)")
		}
	case '\x44':
		if aux == 2 {
			s.Write([]byte{buf})
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
			id[i] = buf[0]
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

	// Debug
	log.Println(buf)

	if buf[0] == '\x4f' {
		return nil
	} else {
		s.Flush()
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
	buf := make([]byte, 1)
	rcv := make([]byte, 6)

	for i := 0; i < 6; {
		n, err := s.Read(buf)
		if err != nil {
			return err
		}
		rcv[i] = buf[0]
		i = i + n
	}

	if bytes.Equal(buf, footer) {
		return nil
	} else {
		s.Flush()
		return errors.New("footer does not match")
	}
}

func sbs32VersionQuery(s *serial.Port, r *bufio.Reader) error {
	// Get footer byte
	buf, err := r.ReadByte()
	if err != nil {
		return err
	}

	if buf != '\x4F' {
		return errors.New("invalid footer received")
	}

	// Send version data
	data := []byte{
		'\x4F', '\x32', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30',
		'\x30', '\x43', '\x43', '\x55', '\x20', '\x66', '\x6F', '\x72', '\x20', '\x43', '\x50', '\x41', '\x20', '\x53', '\x82', '\x72',
		'\x69', '\x65', '\x20', '\x34', '\x30', '\x30', '\x30', '\x0A', '\x50', '\x47', '\x4D', '\x20', '\x0A', '\x50', '\x47', '\x4D',
		'\x20', '\x76', '\x65', '\x72', '\x73', '\x69', '\x6F', '\x6E', '\x3A', '\x20', '\x32', '\x2E', '\x30'}

	_, err = s.Write(data)
	if err != nil {
		return err
	}
	return nil
}
