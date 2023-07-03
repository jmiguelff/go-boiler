package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grid-x/modbus"
)

func main() {
	handler := modbus.NewTCPClientHandler("192.168.1.68:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveID = 0x01
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	//res, err := client.ReadDiscreteInputs(1, 2)
	res, err := client.ReadInputRegisters(3001, 6)
	// res, err := client.ReadHoldingRegisters(3001, 6)
	fmt.Println(res)
}
