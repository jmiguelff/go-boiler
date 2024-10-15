package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	dstPtr := flag.String("dst", "localhost:8000", "data destination")
	flag.Parse()

	// Get the IP struct from hostname
	tcpAddr, err := net.ResolveTCPAddr("tcp", *dstPtr)
	if err != nil {
		log.Fatalln("Resolve address fail ", err)
	}

	// Connect to the server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln("TCP connection fail ", err)
	}
	// Close socket after return from main
	defer conn.Close()

	for {
		buffer := make([]byte, 1400)
		dataSize, err := conn.Read(buffer)
		if err != nil {
			log.Println("Connection closed")
			return
		}

		// Slice just the part of the slice with data
		data := buffer[:dataSize]
		log.Println("Number of bytes received: ", dataSize)
		//log.Println("Received message: ", string(data))
		log.Printf("Received bytes:\n [%x]\n", data)
	}
}
