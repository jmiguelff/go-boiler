package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// User input for destination address (e.g., ":8000" or "192.168.1.100:3000")
	dstPtr := flag.String("dst", "localhost:8000", "IP and port for the server to listen on")
	flag.Parse()

	// Validate the user input
	host, port, err := net.SplitHostPort(*dstPtr)
	if err != nil {
		log.Fatalf("Invalid address: %v", err)
	}

	fmt.Printf("Server listening on %s:%s\n", host, port)
	listener, err := net.Listen("tcp", *dstPtr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("New connection!")
		go listenConnection(conn)
	}
}

func listenConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", clientAddr)

	for {
		buffer := make([]byte, 1400)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Connection closed by client %s: %v\n", clientAddr, err)
			return
		}
		data := buffer[:n]
		msg := string(data)
		fmt.Printf("Received message: %q\n", msg)

		// Reply with Pong
		var reply string
		if msg == "ping" {
			reply = "pong"
		} else {
			reply = "unknown"
		}
		_, err = conn.Write([]byte(reply))
		if err != nil {
			log.Printf("Failed to send reply to %s: %v\n", clientAddr, err)
			return
		}
		fmt.Printf("Sent reply: %q\n", reply)

		// Wait a second (optional on server side)
		time.Sleep(1 * time.Second)
	}
}
