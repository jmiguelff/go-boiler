package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

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
	fd, err := ioutil.ReadFile("settings.yml")
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

	// Start serial port
	cSerial := new(serial.Config)
	cSerial.Name = options.SerialConf.Device
	cSerial.Size = byte(options.SerialConf.Size)
	cSerial.Baud = options.SerialConf.Baud
	cSerial.StopBits = serial.StopBits(options.SerialConf.Stopbits)
	cSerial.Parity = serial.Parity(options.SerialConf.Parity[0])
	cSerial.ReadTimeout = time.Millisecond * time.Duration(options.SerialConf.Timeout)

	fmt.Println("Opening serial port")
	/*
		sfd, err := serial.OpenPort(cSerial)
		if err != nil {
			log.Fatalln(err)
		}
		// Close serial port after usage
		defer sfd.Close()
	*/

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

	// For debug test publish data
	publish(mqttCli, "Ola", "topic/test")

	mqttCli.Disconnect(250)
	// Read data from booster
	/*
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

			// str := hex.EncodeToString(msg)
			str := string(msg[:])
			log.Println(str)
		}
	*/
}

/*
func generateChecksum(msg []byte) byte {
	var sum byte
	for _, b := range msg {
		sum = sum + b
	}
	return ^sum
}
*/

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
