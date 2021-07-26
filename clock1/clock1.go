package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	portPtr := flag.String("port", "8000", "server port")
	flag.Parse()

	dest := "localhost:" + *portPtr
	listener, err := net.Listen("tcp", dest)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
