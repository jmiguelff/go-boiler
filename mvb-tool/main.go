package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

type options struct {
	Client clientOpts `json:"client"`
	Server serverOpts `json:"server"`
}

type clientOpts struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type serverOpts struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func main() {
	optsPtr := flag.String("l", "localhost:8000", "server:port")
	flag.Parse()

	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	var s *net.UDPAddr
	var err error
	var srvName string

	if flagset["l"] {
		fmt.Printf("Server set via flags\n")
		srvName = *optsPtr
		s, err = net.ResolveUDPAddr("udp4", srvName)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Printf("Server not explicitly set, using json confs\n")
		fd, err := os.ReadFile("settings.json")
		if err != nil {
			log.Fatalln(err)
		}

		// Unmarshal json file
		var opts options
		if err = json.Unmarshal(fd, &opts); err != nil {
			log.Fatalln(err)
		}

		srvName := opts.Server.IP + ":" + strconv.Itoa(opts.Server.Port)
		s, err = net.ResolveUDPAddr("udp4", srvName)
		if err != nil {
			log.Fatalln(err)
		}
	}

	conn, err := net.ListenUDP("udp", s)
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()
	log.Printf("Server listening %s\n", conn.LocalAddr().String())
	fmt.Println("")

	for {
		message := make([]byte, 1024)
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		// Binary
		log.Printf("Received data from [%s]\n", remote)
		log.Printf("Received Bytes:\n [%x] (38)\n", message[:rlen])

		// Parsed data
		mvbSize := binary.BigEndian.Uint16(message[0:2])
		log.Printf("MVB port size [%d]\n", mvbSize)

		mvbPort := binary.BigEndian.Uint16(message[2:4])
		log.Printf("MVB port indentifier [%d]\n", mvbPort)

		mvbFresh := binary.BigEndian.Uint16(message[36:38])
		log.Printf("MVB port freshness [%d]\n", mvbFresh)

		if mvbFresh != 65535 {
			log.Printf("\t\tMVB port [%d] may have data\n", mvbPort)
		}

		fmt.Println("")
	}
}
