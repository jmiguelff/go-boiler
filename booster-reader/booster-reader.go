package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tarm/serial"
	"gopkg.in/yaml.v2"
)

// user_options_t struct used to keep the YAML data
type user_options_t struct {
	SerialConf struct {
		Name     string `yaml:"name"`
		Device   string `yaml:"device"`
		Size     int    `yaml:"dataBits"`
		Baud     int    `yaml:"baud"`
		Stopbits int    `yaml:"stopbits"`
		Parity   string `yaml:"parity"`
		Timeout  int    `yaml:"timeout"`
	} `yaml:"serialIn"`
	MqttConf struct {
		ClientId string `yaml:"clientId"`
		Broker   string `yaml:"broker"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Pw       string `yaml:"pw"`
		Topic    string `yaml:"topic"`
	} `yaml:"mqttOut"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v", err)
}

func main() {
	// Open yaml file
	fd, err := os.ReadFile("settings.yml")
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshal yaml file
	var options user_options_t
	err = yaml.Unmarshal(fd, &options)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Serial port name: %s\n", options.SerialConf.Name)
	fmt.Printf("\tDevice: %s\n", options.SerialConf.Device)
	fmt.Printf("\tDatabits: %d\n", options.SerialConf.Size)
	fmt.Printf("\tBaudrate: %d\n", options.SerialConf.Baud)
	fmt.Printf("\tStopbits: %d\n", options.SerialConf.Stopbits)
	fmt.Printf("\tParity: %s\n", options.SerialConf.Parity)
	fmt.Printf("\tTimeout: %d\n", options.SerialConf.Timeout)

	fmt.Printf("MQTT client id: %s\n", options.MqttConf.ClientId)
	fmt.Printf("\tBroker: %s:%d\n", options.MqttConf.Broker, options.MqttConf.Port)
	fmt.Printf("\tCredentials: %s - %s\n", options.MqttConf.User, options.MqttConf.Pw)
	fmt.Printf("\tTopic: %s\n", options.MqttConf.Topic)

	// Start serial port
	cSerial := new(serial.Config)
	cSerial.Name = options.SerialConf.Device
	cSerial.Size = byte(options.SerialConf.Size)
	cSerial.Baud = options.SerialConf.Baud
	cSerial.StopBits = serial.StopBits(options.SerialConf.Stopbits)
	cSerial.Parity = serial.Parity(options.SerialConf.Parity[0])
	cSerial.ReadTimeout = time.Millisecond * time.Duration(options.SerialConf.Timeout)

	fmt.Println("Configuring serial port")
	err = set_uart_to_485h(options.SerialConf.Device)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Opening serial port")
	sfd, err := serial.OpenPort(cSerial)
	if err != nil {
		log.Fatalln(err)
	}
	// Close serial port after usage
	defer sfd.Close()

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(fmt.Sprintf("tcp://%s:%d", options.MqttConf.Broker, options.MqttConf.Port))
	mqttOpts.SetClientID(options.MqttConf.ClientId)
	mqttOpts.SetUsername(options.MqttConf.User)
	mqttOpts.SetPassword(options.MqttConf.Pw)
	mqttOpts.SetDefaultPublishHandler(messagePubHandler)
	mqttOpts.OnConnect = connectHandler
	mqttOpts.OnConnectionLost = connectLostHandler

	mqttCli := mqtt.NewClient(mqttOpts)
	if token := mqttCli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

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
			log.Printf("Invalid message length [%d]\n", msgLen)
			continue
		}

		// Check SOF
		if msg[0] != '\x02' {
			log.Printf("Invalid SOF [%X]\n", msg[0])
			continue
		}

		// str := hex.EncodeToString(msg)
		str := string(string(msg[1 : len(msg)-2]))
		now := time.Now()
		timestamp := now.Format("2006-01-02 15:04:05")

		message := fmt.Sprintf("[%s] %s", timestamp, str)
		fmt.Println(message)
		publish(mqttCli, message, options.MqttConf.Topic)
	}
}

/*
func generateChecksum(msg []byte) byte {
	var sum byte
	for _, b := range msg {
		sum = sum + b
	}
	return ^sum
}*/

func publish(client mqtt.Client, msg string, topic string) {
	token := client.Publish(topic, 1, false, msg)
	token.Wait()
}

/*
func sub(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}*/

func set_uart_to_485h(device string) error {
	// change gpio pin attribute
	uartModePath, err := get_uart_mode_path(device)
	if err != nil {
		return err
	}

	uartDuplexPath, err := get_uart_duplex_path(device)
	if err != nil {
		return err
	}

	uartModeFile, err := os.Create(uartModePath)
	if err != nil {
		return err
	}
	defer uartModeFile.Close()

	_, err = io.WriteString(uartModeFile, "1")
	if err != nil {
		return err
	}

	uartDuplexFile, err := os.Create(uartDuplexPath)
	if err != nil {
		return err
	}
	defer uartDuplexFile.Close()

	_, err = io.WriteString(uartDuplexFile, "1")
	if err != nil {
		return err
	}

	// IOCTL Call
	fd, err := syscall.Open(device, syscall.O_RDWR|syscall.O_NDELAY|syscall.O_NOCTTY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	var rs485Flag uint32 = 0x00000003
	const TIOCSRS485 int = 0x542F

	// IOCTL to TIOCSRS485 0x542F
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(TIOCSRS485), uintptr(unsafe.Pointer(&rs485Flag)))
	if errno != 0 {
		fmt.Println("ioctl failed:", errno)
		return fmt.Errorf("ioctl failed [%v]", errno)
	}

	return nil
}

func get_uart_mode_path(device string) (string, error) {
	var uartModePath string

	switch device {
	case "/dev/ttymxc0":
		uartModePath = "/var/gpio/UART1_MODE/value"
	case "/dev/ttymxc1":
		uartModePath = "/var/gpio/UART2_MODE/value"
	case "/dev/ttymxc2":
		uartModePath = "/var/gpio/UART3_MODE/value"
	case "/dev/ttymxc3":
		uartModePath = "/var/gpio/UART4_MODE/value"
	default:
		return uartModePath, fmt.Errorf("invalid device [%s]", device)
	}

	return uartModePath, nil
}

func get_uart_duplex_path(device string) (string, error) {
	var uartDuplexPath string

	switch device {
	case "/dev/ttymxc0":
		uartDuplexPath = "/var/gpio/UART1_HDPLX/value"
	case "/dev/ttymxc1":
		uartDuplexPath = "/var/gpio/UART2_HDPLX/value"
	case "/dev/ttymxc2":
		uartDuplexPath = "/var/gpio/UART3_HDPLX/value"
	case "/dev/ttymxc3":
		uartDuplexPath = "/var/gpio/UART4_HDPLX/value"
	default:
		return uartDuplexPath, fmt.Errorf("invalid device [%s]", device)
	}

	return uartDuplexPath, nil
}
