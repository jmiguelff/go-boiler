package main

import (
	"flag"
	"log"
	"net"
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
	log.Printf("server listening %s\n", conn.LocalAddr().String())
	log.Printf("server listening %s:%d\n", *address, *port)

	for {
		message := make([]byte, 1024)
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		// Binary
		log.Printf("received-bytes: [%x] from [%s]\n", message[:rlen], remote)

		// String
		//data := strings.TrimSpace(string(message[:rlen]))
		//log.Printf("received: %s from %s\n", data, remote)
	}
}
