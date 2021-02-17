package main

import (
	"bufio"
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

	// Get EM acknowledge (q -> "/?!<CR><LF>")
	fmt.Println("Try to get acknowledge from Emeter")
	q := []byte{'/', '?', '!', '\x0d', '\x0a'}
	_, err = sfd.Write(q)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(sfd)
	reply, err := reader.ReadBytes('\x0a')
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(string(reply))

	// Read from table 1
	//fmt.Println("Try to read from table 1")
	//q1 := []byte{0x06, '0', '5', '0', 0x0d, 0x0a}
	//_, err = sfd.Write(q1)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//n, err = sfd.Read(r)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(string(r[:n]))

	// Exit
	e := []byte{'\x01', 'B', '0', '\x03', 'q'}
	_, err = sfd.Write(e)
	if err != nil {
		log.Fatal(err)
	}

	// Exit
	fmt.Println("Exit")
	return
}
