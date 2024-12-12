package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/tarm/serial"
)

func main() {
	dstPtr := flag.String("dst", "localhost:8000", "data destination")
	devPtr := flag.String("dev", "/dev/ttymxc4", "serial device")
	baudPtr := flag.Int("baud", 9600, "serial baud")
	flag.Parse()

	c := &serial.Config{Name: *devPtr, Baud: *baudPtr}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}

	listener, err := net.Listen("tcp", *dstPtr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", *dstPtr, err)
	}
	defer listener.Close()

	log.Printf("Listening on %s for TCP connections", *dstPtr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		log.Printf("New connection from %s", conn.RemoteAddr())

		// Start bidirectional copying
		go handleConnection(conn, s)
	}
}

func handleConnection(conn net.Conn, s *serial.Port) {
	defer conn.Close()

	// Use channels to signal completion of copy operations.
	done := make(chan struct{})

	// Copy serial->TCP
	go func() {
		_, err := io.Copy(conn, s)
		if err != nil {
			log.Printf("Error copying from serial to TCP: %v", err)
		}
		done <- struct{}{}
	}()

	// Copy TCP->serial
	go func() {
		_, err := io.Copy(s, conn)
		if err != nil {
			log.Printf("Error copying from TCP to serial: %v", err)
		}
		done <- struct{}{}
	}()

	// Wait for either direction to end
	<-done
	// When one direction ends (e.g., client disconnects), close the connection.
	conn.Close()

	// The other goroutine will return from io.Copy due to the closed connection.
	<-done
	log.Printf("Connection from %s closed", conn.RemoteAddr())
}
