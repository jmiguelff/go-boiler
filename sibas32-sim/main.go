package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/tarm/serial"
)

type state int

const (
	WAITFORHEADER state = iota
	WAITFORCMD
	REQCONNECTORBYID
	REQBAUDRATECHANGE
	REQSWVERSION
	REQEVENTS
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

	// initial conditions
	header := 0
	cmdOk := 0
	selCmd := NA
	state := WAITFORHEADER

	// infinite loop
	for {
		switch state {
		case WAITFORHEADER:
			if header == 1 {
				state = WAITFORCMD
			} else if header == 2 {
				state = REQEVENTS
			} else {
				state = WAITFORHEADER
			}
		case WAITFORCMD:
			if selCmd == SETBAUDRATE {
				state = REQBAUDRATECHANGE
			} else if selCmd == GETCONNECTORBYID {
				state = REQCONNECTORBYID
			} else if selCmd == GETSWVERSION {
				state = REQSWVERSION
			} else if selCmd == NA && header == 1 {
				state = WAITFORCMD
			} else if selCmd == NA && header == 0 {
				state = WAITFORHEADER
			} else {
				state = ERROR
			}
		case REQCONNECTORBYID:
			if cmdOk == 0 {
				state = REQCONNECTORBYID
			} else {
				state = WAITFORHEADER
			}
		case REQBAUDRATECHANGE:
			if cmdOk == 0 {
				state = REQBAUDRATECHANGE
			} else {
				state = WAITFORHEADER
			}
		case REQSWVERSION:
			if cmdOk == 0 {
				state = REQSWVERSION
			} else {
				state = WAITFORHEADER
			}
		case REQEVENTS:
			if cmdOk == 0 {
				state = REQEVENTS
			} else {
				state = WAITFORHEADER
			}
		default:
			state = ERROR
		}

		switch state {
		case WAITFORHEADER:
			cmdOk = 0
			selCmd = NA
			header, err = sbs32WaitForHeader(sfd, r)
			if err != nil {
				log.Println(err)
			}
		case WAITFORCMD:
			header = 0
			selCmd, err = sbs32WaitForCmd(sfd, r)
			if err != nil {
				selCmd = NA
				log.Println(err)
			}
		case REQCONNECTORBYID:
			err = sbs32ReqConnByIdResp(sfd, r)
			if err != nil {
				cmdOk = -1
				log.Println(err)
				break
			}
			cmdOk = 1
		case REQBAUDRATECHANGE:
			err = sbs32ReqBaudChangeResp(sfd, r, c)
			if err != nil {
				cmdOk = -1
				log.Println(err)
				break
			}
			cmdOk = 1
		case REQSWVERSION:
			err = sbs32VersionResp(sfd, r)
			if err != nil {
				cmdOk = -1
				log.Println(err)
				break
			}
			cmdOk = 1
		case REQEVENTS:
			header = 0
			err = sbs32ReqEventsResp(sfd, r)
			if err != nil {
				cmdOk = -1
				log.Println(err)
				break
			}
			cmdOk = 1
		case ERROR:
			log.Fatalln("Fatal error - exit")
		default:
			log.Fatalln("Failed to default state on actions switch/case")
		}
	}
}

func sbs32WaitForHeader(s *serial.Port, r *bufio.Reader) (int, error) {
	// wait for first header byte
	log.Println("sbs32WaitForHeader: Waiting for header")
	for {
		buf, err := r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("sbs32WaitForHeader: %w", err)
		}
		if buf == '\x00' {
			// Debug
			log.Printf("Received first header byte [%x]\n", buf)
			break
		}
	}

	// get second header byte
	buf, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("sbs32WaitForHeader: %w", err)
	}

	// If second byte
	if buf != '\xF1' && buf != '\x4F' {
		return 0, errors.New("sbs32WaitForHeader: second header byte does not match protocol")
	}

	// Debug
	log.Printf("Received second header byte [%x]\n", buf)

	if buf == '\xF1' {
		// send ack byte for standard header
		_, err = s.Write([]byte{'\xF2'})
		if err != nil {
			return 0, fmt.Errorf("sbs32WaitForHeader: %w", err)
		}

		// get ack echo byte
		buf, err = r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("sbs32WaitForHeader: %w", err)
		}

		if buf != '\xF2' {
			return 0, errors.New("ack echo byte does not match")
		}

		return 1, nil

	} else {
		// send echo for events header
		_, err = s.Write([]byte{'\x4F'})
		if err != nil {
			return 0, fmt.Errorf("sbs32WaitForHeader: %w", err)
		}

		return 2, nil
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

func sbs32WaitForFPId(s *serial.Port, r *bufio.Reader) ([]byte, error) {
	id := make([]byte, 8)

	for i := 0; i < len(id); i++ {
		buf, err := r.ReadByte()
		if err != nil {
			return id, fmt.Errorf("sbs32WaitForFPId: %w", err)
		} else {
			s.Write([]byte{buf})
			id[i] = buf
		}
	}

	return id, nil
}

func sbs32WaitForFooter(r *bufio.Reader) error {
	buf, err := r.ReadByte()
	if err != nil {
		return fmt.Errorf("sbs32WaitForFooter: %w", err)
	}

	if buf != '\x4F' {
		return errors.New("sbs32WaitForFooter: invalid footer received")
	}

	return nil
}

func sbs32SendConnData(s *serial.Port) error {
	data := []byte{
		'\x45', '\x4E', '\x41', '\x42', '\x44', '\x49', '\x41', '\x33', '\x2E', '\x4F', '\x55', '\x54', '\x33', '\x36', '\x30', '\x30',
		'\x30', '\x30', '\x30', '\x30', '\x4F', '\x39', '\x45', '\x30', '\x30', '\x45', '\x42', '\x32', '\x43', '\x30', '\x30', '\x32',
		'\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x38', '\x46', '\x33', '\x20', '\x20', '\x20', '\x20', '\x20', '\x20', '\x20',
		'\x20', '\x20', '\x20', '\x32', '\x36', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x31', '\x30', '\x30', '\x30', '\x30',
		'\x30', '\x30', '\x30'}

	_, err := s.Write(data)
	if err != nil {
		return fmt.Errorf("sbs32SendConnData: %w", err)
	}

	return nil
}

func sbs32EchoPayloadBaud(s *serial.Port, r *bufio.Reader) error {
	// Receive 4 bytes of payload and just echo back
	for i := 0; i < 4; i++ {
		buf, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("sbs32EchoPayloadBaud: %w", err)
		}
		_, err = s.Write([]byte{buf})
		if err != nil {
			return fmt.Errorf("sbs32EchoPayloadBaud: %w", err)
		}
	}

	return nil
}

func sbs32WaitForFooterBaud(r *bufio.Reader) error {
	footer := []byte{'\x4F', '\xFC', '\xFC', '\xFC', '\xFC', '\xFC'}

	// Receive 6 bytes
	for i := 0; i < 6; {
		buf, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("sbs32WaitForFooterBaud: %w", err)
		}
		if buf != footer[i] {
			return errors.New("sbs32WaitForFooterBaud: footer does not match")
		}
	}

	return nil
}

func sbs32VersionResp(s *serial.Port, r *bufio.Reader) error {
	// Get footer byte
	err := sbs32WaitForFooter(r)
	if err != nil {
		return fmt.Errorf("sbs32VersionResp: %w", err)
	}

	// Send version data
	data := []byte{
		'\x4F', '\x32', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30', '\x30',
		'\x30', '\x43', '\x43', '\x55', '\x20', '\x66', '\x6F', '\x72', '\x20', '\x43', '\x50', '\x41', '\x20', '\x53', '\x82', '\x72',
		'\x69', '\x65', '\x20', '\x34', '\x30', '\x30', '\x30', '\x0A', '\x50', '\x47', '\x4D', '\x20', '\x0A', '\x50', '\x47', '\x4D',
		'\x20', '\x76', '\x65', '\x72', '\x73', '\x69', '\x6F', '\x6E', '\x3A', '\x20', '\x32', '\x2E', '\x30'}

	_, err = s.Write(data)
	if err != nil {
		return fmt.Errorf("sbs32VersionResp: %w", err)
	}

	return nil
}

func sbs32ReqConnByIdResp(s *serial.Port, r *bufio.Reader) error {
	// Wait for FP ID
	adr, err := sbs32WaitForFPId(s, r)
	if err != nil {
		return fmt.Errorf("sbs32ReqConnByIdResp: %w", err)
	}

	// Debug
	log.Println(adr)

	err = sbs32WaitForFooter(r)
	if err != nil {
		return fmt.Errorf("sbs32ReqConnByIdResp: %w", err)
	}

	err = sbs32SendConnData(s)
	if err != nil {
		return fmt.Errorf("sbs32ReqConnByIdResp: %w", err)
	}

	return nil
}

func sbs32ReqBaudChangeResp(s *serial.Port, r *bufio.Reader, c *serial.Config) error {
	err := sbs32EchoPayloadBaud(s, r)
	if err != nil {
		return fmt.Errorf("sbs32ReqBaudChangeResp: %w", err)
	}

	err = sbs32WaitForFooterBaud(r)
	if err != nil {
		return fmt.Errorf("sbs32ReqBaudChangeResp: %w", err)
	}

	// Change serial port baudrate to 115200
	s.Close()
	c.Baud = 115200
	s, err = serial.OpenPort(c)
	if err != nil {
		return fmt.Errorf("sbs32ReqBaudChangeResp: %w", err)
	}

	return nil
}

func sbs32ReqEventsResp(s *serial.Port, r *bufio.Reader) error {
	// Receive and echo 21 bytes
	// This byte have the type of events requested
	for i := 0; i < 21; i++ {
		_, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("sbs32ReqEventsResp: %w", err)
		}

		buf, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("sbs32ReqEventsResp: %w", err)
		}

		_, err = s.Write([]byte{buf})
		if err != nil {
			return fmt.Errorf("sbs32ReqEventsResp: %w", err)
		}
	}

	// Get cmd footer
	buf, err := r.ReadByte()
	if err != nil {
		return fmt.Errorf("sbs32ReqEventsResp: %w", err)
	}

	if buf != '\x4E' {
		return errors.New("sbs32ReqEventsResp: command footer 1st byte does not match")
	}

	buf, err = r.ReadByte()
	if err != nil {
		return fmt.Errorf("sbs32ReqEventsResp: %w", err)
	}

	if buf != '\x4F' {
		return errors.New("sbs32ReqEventsResp: command footer 2nd byte does not match")
	}

	// Get data to send
	fd, err := os.Open("data/cpa4000-events-example.txt")
	if err != nil {
		return fmt.Errorf("sbs32ReqEventsResp: %w", err)
	}
	defer fd.Close()
	fr := bufio.NewReader(fd)

	for {
		d, err := fr.ReadByte()
		if err != nil {
			log.Printf("sbs32ReqEventsResp: %s\n", err)
			break
		}

		// Send data byte
		_, err = s.Write([]byte{d})
		if err != nil {
			return fmt.Errorf("sbs32ReqEventsResp: %w", err)
		}

		// Receive echo byte
		buf, err = r.ReadByte()
		if err != nil {
			return fmt.Errorf("sbs32ReqEventsResp: %w", err)
		}

		if buf == d {
			continue
		}

		if buf == '\x23' {
			log.Println("Received stop command")
			break
		}
	}

	// End of frame
	_, err = s.Write([]byte{'\x02'})
	if err != nil {
		return fmt.Errorf("sbs32ReqEventsResp: %w", err)
	}

	return nil
}
