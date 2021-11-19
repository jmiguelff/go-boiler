package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type metadata struct {
	DataType    string `json:"type"`
	Name        string `json:"id"`
	DataSize    int    `json:"size"`
	DeviceId    string `json:"deviceId"`
	DataCounter int    `json:"counter"`
	TimeStamp   string `json:"time"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if msg.Topic() == "cp/obd/uqe2400/2400/m1/mb1000/datalogger/signals" {
		h := bytes.Split(msg.Payload(), []byte("\n"))
		res := metadata{}
		if err := json.Unmarshal(h[0], &res); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Received file [%s] on topic [%s]\n", res.Name, msg.Topic())
		path := "data/" + res.Name
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		for _, line := range h {
			line = append(line, '\n')
			_, err = file.Write(line)
			if err != nil {
				panic(err)
			}
		}
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	var broker = "192.168.3.53"
	var port = 1883

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("mqtt-sub-cli")
	opts.SetUsername("test")
	opts.SetPassword("n0madr00t")
	opts.SetDefaultPublishHandler(messagePubHandler)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)

	<-c
}

func sub(client mqtt.Client) {
	topic := "cp/obd/uqe2400/2400/m1/mb1000/datalogger/signals"
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic %s\n", topic)
}
