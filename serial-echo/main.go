package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
	"gopkg.in/yaml.v2"
)

type serialConfigT struct {
	SerialMode struct {
		Name     string `yaml:"name"`
		Device   string `yaml:"device"`
		DataSize int    `yaml:"dataSize"`
		Baud     int    `yaml:"baud"`
		Stopbits int    `yaml:"stopbits"`
		Parity   string `yaml:"parity"`
		Timeout  int    `yaml:"timeout"`
	} `yaml:"serial"`
}

func setSerialMode(c *serialConfigT) *serial.Config {
	s := new(serial.Config)
	s.Name = c.SerialMode.Device
	s.Baud = c.SerialMode.Baud
	s.Size = byte(c.SerialMode.DataSize)
	s.StopBits = serial.StopBits(c.SerialMode.Stopbits)
	s.Parity = serial.Parity(c.SerialMode.Parity[0])
	s.ReadTimeout = time.Millisecond * time.Duration(c.SerialMode.Timeout)

	return s
}

func main() {
	// Open settings file
	fd, err := ioutil.ReadFile("settings.yml")
	if err != nil {
		log.Fatalln(err)
	}

	// Parse settings file (YAML)
	opts := new(serialConfigT)
	err = yaml.Unmarshal(fd, opts)
	if err != nil {
		log.Fatalln(err)
	}

	// Basic gui
	log.Println("Test serial port app")
	log.Println("Serial port configurations")
	log.Println(*opts)

	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Press enter to start")
	scanner.Scan()

	// Serial port configuration
	mode := setSerialMode(opts)

	log.Println("Open serial port")
	sfd, err := serial.OpenPort(mode)
	if err != nil {
		log.Fatalln(err)
	}
	defer sfd.Close()

	testVals := []byte{'n', 'o', 'm', 'a', 'd', ' ', 't', 'e', 'c', 'h'}
	for _, val := range testVals {
		if sendByteWithEcho(sfd, val) < 0 {
			log.Fatalln("Serial port test [fail]")
			return
		}
	}

	log.Println("Serial port test [pass]")
	return
}

func sendByteWithEcho(sfd *serial.Port, b byte) int {
	log.Printf("TX -> [%c]\n", b)
	_, err := sfd.Write([]byte{b})
	if err != nil {
		return -1
	}

	r := bufio.NewReader(sfd)
	ret, err := r.ReadByte()
	if err != nil {
		return -1
	}
	log.Printf("RX -> [%c]\n", ret)

	if ret != b {
		log.Fatalf("Echo to byte [%c] does not match [%c]\n", b, ret)
		return -1
	}

	return 0
}
