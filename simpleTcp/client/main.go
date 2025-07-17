package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	dstPtr := flag.String("dst", "localhost:8000", "server address (host:port)")
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", *dstPtr)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Fail to connect to %s: %v", *dstPtr, err)
	}
	defer conn.Close()
	fmt.Printf("Connected to %s\n", *dstPtr)

	for {
		// Send ping
		_, err := conn.Write([]byte("ping"))
		if err != nil {
			log.Printf("Disconnected (write error): %v", err)
			return
		}
		fmt.Println("Sent: ping")

		// Read pong
		replyBuf := make([]byte, 1400)
		n, err := conn.Read(replyBuf)
		if err != nil {
			log.Printf("Disconnected (read error): %v", err)
			return
		}
		fmt.Printf("Received: %q\n", string(replyBuf[:n]))

		time.Sleep(1 * time.Second)
	}
}
