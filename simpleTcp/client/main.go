package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// TODO: Change client parameters and message to send user input
	// Get the IP struct from hostname
	servAddr := "localhost:3000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to the server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Fail to connect to ", servAddr)
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

	// Handle the retry from server
	reply := make([]byte, 1400)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Reply from server: ", string(reply))

}
