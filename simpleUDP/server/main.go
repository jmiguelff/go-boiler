package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {
	address := flag.String("a", "localhost", "server address")
	port := flag.Int("p", 8000, "server port")
	flag.Parse()

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: *port,
		IP:   net.ParseIP(*address),
	})
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())
	fmt.Printf("server listening %s:%d\n", *address, *port)

	for {
		message := make([]byte, 1024)
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		// Binary
		fmt.Printf("received-bytes: [%x]\n", message[:rlen])

		// String
		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Printf("received: %s from %s\n", data, remote)
	}
}
