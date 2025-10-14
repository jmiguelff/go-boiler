package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
	"gopkg.in/yaml.v2"
)

// SerialOptT holds the YAML config.
type SerialOptT struct {
	SerialConf struct {
		Device   string `yaml:"device"`
		Size     int    `yaml:"dataBits"`
		Baud     int    `yaml:"baud"`
		Stopbits int    `yaml:"stopbits"`
		Parity   string `yaml:"parity"`
		Timeout  int    `yaml:"timeout"`
	} `yaml:"serialConf"`
}

func main() {
	// 1) Load YAML config
	data, err := os.ReadFile("settings.yml")
	if err != nil {
		log.Fatalf("read settings.yml: %v", err)
	}
	var cfg SerialOptT
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("parse YAML: %v", err)
	}

	// 2) Configure & open serial port
	portCfg := &serial.Config{
		Name:        cfg.SerialConf.Device,
		Baud:        cfg.SerialConf.Baud,
		Size:        byte(cfg.SerialConf.Size),
		StopBits:    serial.StopBits(cfg.SerialConf.Stopbits),
		Parity:      serial.Parity(cfg.SerialConf.Parity[0]),
		ReadTimeout: time.Millisecond * time.Duration(cfg.SerialConf.Timeout),
	}

	log.Printf("Opening serial port %s @ %d baud", portCfg.Name, portCfg.Baud)
	sfd, err := serial.OpenPort(portCfg)
	if err != nil {
		log.Fatalf("open port: %v", err)
	}
	defer sfd.Close()

	log.Println("Serial responder started. Waiting for '0x0A' to respond with 'WORKING'...")

	// 3) Read from serial port and respond to '\n'
	buf := make([]byte, 128)
	for {
		n, err := sfd.Read(buf)
		if err != nil {
			// Timeout is not an error, just continue
			continue
		}

		if n > 0 {
			// Process received data
			for i := 0; i < n; i++ {
				char := buf[i]

				// Echo to stdout
				fmt.Printf("Received: %c (0x%02X)\n", char, char)

				// Check if it's '\n'
				if char == '\n' {
					response := "WORKING"
					_, err := sfd.Write([]byte(response))
					if err != nil {
						log.Printf("Error writing response: %v", err)
					} else {
						log.Printf("Sent response: %s", response)
					}
				}
			}
		}
	}
}
