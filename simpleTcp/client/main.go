package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	dstPtr := flag.String("dst", "localhost:8000", "data destinantion")
	flag.Parse()

	// Get the IP struct from hostname
	tcpAddr, err := net.ResolveTCPAddr("tcp", *dstPtr)
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to the server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Fail to connect to ", *dstPtr)
		log.Fatalln(err)
	}
	// Close socket after return from main
	defer conn.Close()

	// Send data
	strEcho := "Test message"
	_, err = conn.Write([]byte(strEcho))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Message sent: ", strEcho)

	// Handle the reply from server
	reply := make([]byte, 1400)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Reply from server: ", string(reply))

}
