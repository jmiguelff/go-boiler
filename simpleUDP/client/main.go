package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	// Define destination and command flags
	dstPtr := flag.String("dst", "localhost:8000", "data destination")
	cmdPtr := flag.String("cmd", "Test message", "command to send")
	stopStringPtr := flag.String("stop", "END", "string to stop receiving data")
	srcPtr := flag.String("cli", "localhost:8000", "client src") // Add client port flag
	flag.Parse()

	// Resolve the address (for UDP)
	clientUDPAddr, err := net.ResolveUDPAddr("udp", *dstPtr)
	if err != nil {
		log.Fatalf("Invalid server address: %v\n", err)
	}

	// Validate the src address using net.SplitHostPort
	myIP, myPort, err := net.SplitHostPort(*srcPtr)
	if err != nil {
		log.Fatalf("Invalid client address: %v", err)
	}

	// Convert the port from string to int
	portInt, err := strconv.Atoi(myPort)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	// Parse the IP address
	ip := net.ParseIP(myIP)
	if ip == nil {
		log.Fatalf("Invalid IP address: %s", myIP)
	}

	// Listen on the specified UDP address and port
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: portInt,
		IP:   ip,
	})
	if err != nil {
		log.Fatalf("Failed to set up UDP server: %v\n", err)
	}
	defer conn.Close()

	// Send user-defined command
	command := *cmdPtr + "\n"
	_, err = conn.WriteToUDP([]byte(command), clientUDPAddr)
	if err != nil {
		log.Fatalf("Failed to send data: %v\n", err)
	}
	fmt.Printf("Command sent: %s\n", command)

	// Start receiving data until the stop string is found
	reply := make([]byte, 1400)
	receivedData := ""

	for {
		n, _, err := conn.ReadFromUDP(reply)
		if err != nil {
			log.Fatalf("Failed to receive data: %v\n", err)
		}

		receivedData += string(reply[:n])
		fmt.Println("Received data:", string(reply[:n]))

		// Check if the stop string appears in the received data
		if strings.Contains(receivedData, *stopStringPtr) {
			fmt.Printf("Stop string '%s' found. Stopping reception.\n", *stopStringPtr)
			break
		}
	}
}
