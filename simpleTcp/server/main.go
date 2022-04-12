package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// TODO: Change server parameters to be user inputs
	fmt.Println("Server listening on 3000")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalln(err)
	}
	// Close socket after return from main
	defer listener.Close()

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
	for {
		buffer := make([]byte, 1400)
		dataSize, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Connection closed")
			return
		}

		// Slice just the part of the slice with data
		data := buffer[:dataSize]
		fmt.Println("Received message: ", string(data))
		fmt.Printf("Received bytes: [%x]\n", data)

		// Echo data back
		// _, err = conn.Write(data)
		// if err != nil {
		// 	 log.Fatalln(err)
		// }
		// fmt.Println("Message sent: ", string(data))
	}
}
