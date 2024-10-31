package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	// User input for destination address (e.g., "localhost:8000" or "192.168.1.100:3000")
	dstPtr := flag.String("dst", "localhost:8000", "IP and port for the server to listen on")
	flag.Parse()

	// Validate the user input using net.SplitHostPort
	host, port, err := net.SplitHostPort(*dstPtr)
	if err != nil {
		log.Fatalf("Invalid address: %v", err)
	}

	// Use the provided IP and port for setting up the server
	fmt.Printf("Server listening on %s:%s\n", host, port)
	listener, err := net.Listen("tcp", *dstPtr) // Listening on the user-defined IP and port
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close() // Close the listener when the application exits

	// Infinite loop accepting connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("New connection!")

		// Create another go routine to deal with data
		// This is like a thread
		go listenConnection(conn)
	}
}

func listenConnection(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when the function returns

	// Print the client's IP and port (origin address)
	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", clientAddr)

	for {
		buffer := make([]byte, 1400)
		dataSize, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Connection closed by client: %s\n", clientAddr)
			return
		}

		// Slice just the part of the slice with data
		data := buffer[:dataSize]
		fmt.Println("Received message: ", string(data))
		fmt.Printf("Received bytes: [%x]\n", data)

		// Echo data back (optional)
		// _, err = conn.Write(data)
		// if err != nil {
		//     log.Fatalln(err)
		// }
		// fmt.Println("Message sent: ", string(data))
	}
}
