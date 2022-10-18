package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/tarm/serial"
	"gopkg.in/yaml.v2"
)

// SerialOptT struct used to keep the YAML data
type SerialOptT struct {
	SerialConf struct {
		Name     string `yaml:"name"`
		Device   string `yaml:"device"`
		Size     int    `yaml:"dataBits"`
		Baud     int    `yaml:"baud"`
		Stopbits int    `yaml:"stopbits"`
		Parity   string `yaml:"parity"`
		Timeout  int    `yaml:"timeout"`
	} `yaml:"serialConf"`
}

func main() {
	// Open yaml file
	fd, err := ioutil.ReadFile("settings.yml")
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshal yaml file
	var serialOpts SerialOptT
	err = yaml.Unmarshal(fd, &serialOpts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Serial port name: %s\n", serialOpts.SerialConf.Name)
	fmt.Printf("\tDevice: %s\n", serialOpts.SerialConf.Device)
	fmt.Printf("\tDatabits: %d\n", serialOpts.SerialConf.Size)
	fmt.Printf("\tBaudrate: %d\n", serialOpts.SerialConf.Baud)
	fmt.Printf("\tStopbits: %d\n", serialOpts.SerialConf.Stopbits)
	fmt.Printf("\tParity: %s\n", serialOpts.SerialConf.Parity)
	fmt.Printf("\tTimeout: %d\n", serialOpts.SerialConf.Timeout)

	// Start serial port
	cSerial := new(serial.Config)
	cSerial.Name = serialOpts.SerialConf.Device
	cSerial.Size = byte(serialOpts.SerialConf.Size)
	cSerial.Baud = serialOpts.SerialConf.Baud
	cSerial.StopBits = serial.StopBits(serialOpts.SerialConf.Stopbits)
	cSerial.Parity = serial.Parity(serialOpts.SerialConf.Parity[0])
	cSerial.ReadTimeout = time.Millisecond * time.Duration(serialOpts.SerialConf.Timeout)

	fmt.Println("Opening serial port")

	sfd, err := serial.OpenPort(cSerial)
	if err != nil {
		log.Fatalln(err)
	}

	// Close serial port after usage
	defer sfd.Close()

	// Read data from booster
	for {
		reader := bufio.NewReader(sfd)
		msg, err := reader.ReadBytes('\x0a')
		if err != nil {
			log.Fatalln(err)
		}

		// Check min size
		msgLen := len(msg)
		if msgLen < 8 {
			continue
		}

		// Check SOF
		if msg[0] != '\x02' {
			continue
		}

		// Check payload length
		if int(msg[1]) != (msgLen - 1) {
			continue
		}

		// TODO: Test checksum
		cksum := generateChecksum(msg)
		if cksum != msg[msgLen-2] {
			continue
		}

		str := hex.EncodeToString(msg)
		log.Println(str)
	}
}

func generateChecksum(msg []byte) byte {
	var sum byte
	for _, b := range msg {
		sum = sum + b
	}
	return ^sum
}
