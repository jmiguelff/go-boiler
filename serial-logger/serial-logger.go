package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	// 2) Open logfile & create logger
	f, err := os.OpenFile("serial.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("open serial.log: %v", err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags|log.Lmicroseconds)

	// 3) Configure & open serial port
	portCfg := &serial.Config{
		Name:        cfg.SerialConf.Device,
		Baud:        cfg.SerialConf.Baud,
		Size:        byte(cfg.SerialConf.Size),
		StopBits:    serial.StopBits(cfg.SerialConf.Stopbits),
		Parity:      serial.Parity(cfg.SerialConf.Parity[0]),
		ReadTimeout: time.Millisecond * time.Duration(cfg.SerialConf.Timeout),
	}
	logger.Printf("opening %s @ %d baud", portCfg.Name, portCfg.Baud)
	sfd, err := serial.OpenPort(portCfg)
	if err != nil {
		logger.Fatalf("open port: %v", err)
	}
	defer sfd.Close()

	// 4) Graceful shutdown on CTRL+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Println("interrupt received, shutting down")
		sfd.Close()
		f.Close()
		os.Exit(0)
	}()

	// 5) Frame‑based read & ASCII‑hex log
	reader := bufio.NewReader(sfd)
	var buf []byte

	for {
		b, err := reader.ReadByte()
		if err != nil {
			logger.Printf("read error: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if b == 0x55 {
			// New frame: flush previous buffer
			if len(buf) > 0 {
				// build ASCII‑hex line
				parts := make([]string, len(buf))
				for i, vb := range buf {
					parts[i] = fmt.Sprintf("0x%02X", vb)
				}
				logger.Println(strings.Join(parts, " "))
			}
			buf = []byte{b}
		} else {
			buf = append(buf, b)
		}
	}
}
