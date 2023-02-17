package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type HVAC struct {
	Name string `yaml:"name"`
	Ip   string `yaml:"ip"`
}

func main() {
	address := flag.String("a", "localhost", "server address")
	port := flag.Int("p", 8000, "server port")
	flag.Parse()

	// Open yaml file
	fd, err := os.ReadFile("settings.yaml")
	if err != nil {
		panic(err)
	}

	// Unmarshal yaml file
	var hvacs []HVAC
	err = yaml.Unmarshal(fd, &hvacs)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: *port,
		IP:   net.ParseIP(*address),
	})
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())
	// fmt.Printf("server listening %s:%d\n", *address, *port)

	for {
		message := make([]byte, 1024)
		_, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		for _, v := range hvacs {
			if remote.IP.String() == v.Ip {
				// Version
				fmt.Printf("HVAC [%s] has the following configurations:\n", v.Name)
				fmt.Printf("\tSerial Number [%s]\n", string(message[13:22]))
				fmt.Printf("\tSoftware Version [%s]\n", string(message[23:32]))
				fmt.Printf("\tKernel Version [%s]\n", string(message[33:42]))
			}
		}
	}
}
