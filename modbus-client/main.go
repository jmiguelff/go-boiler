package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grid-x/modbus"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <IP> <Porta>")
		os.Exit(1)
	}

	ip := os.Args[1]
	port := os.Args[2]
	address := fmt.Sprintf("%s:%s", ip, port)

	handler := modbus.NewTCPClientHandler(address)
	handler.Timeout = 10 * time.Second
	handler.SlaveID = 0x01
	handler.Logger = log.New(os.Stdout, "modbus: ", log.LstdFlags)

	// Connect
	err := handler.Connect()
	if err != nil {
		log.Fatalln("Erro ao ligar:", err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	res, err := client.ReadHoldingRegisters(0, 10) // lê 40001 a 40010
	if err != nil {
		log.Fatalln("Erro na leitura:", err)
	}
	fmt.Printf("Resposta: %v\n", res)
}
