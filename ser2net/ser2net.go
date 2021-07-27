package main

import (
	"flag"
	"log"
	"net"

	"github.com/tarm/serial"
)

func main() {
	dstPtr := flag.String("dst", "localhost:8000", "data destinantion")
	devPtr := flag.String("dev", "/dev/ttymxc4", "serial device")
	baudPtr := flag.Int("baud", 9600, "serial baud")
	flag.Parse()

	c := &serial.Config{Name: *devPtr, Baud: *baudPtr}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", *dstPtr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Print("listening on " + *dstPtr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		log.Print("Got a connection")
		go serialToTcp(conn, s)
		go tcpToSerial(conn, s)
	}

}

func serialToTcp(c net.Conn, s *serial.Port) {
	buf := make([]byte, 512)
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		c.Write(buf[:n])
	}
}

func tcpToSerial(c net.Conn, s *serial.Port) {
	buf := make([]byte, 512)
	for {
		n, err := c.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		s.Write(buf[:n])
	}
}
