package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Start dummy server...")

	// Listen to port
	listener, err := net.Listen("tcp", ":4001")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Client connected")

	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buf[:n])
	}
}
