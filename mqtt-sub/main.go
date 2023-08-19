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

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", " "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	h := bytes.Split(msg.Payload(), []byte("\n"))
	res := metadata{}
	if err := json.Unmarshal(h[0], &res); err != nil {
		fmt.Println(err)
	}

	// Print topic
	fmt.Printf("Received data [%s] on topic [%s]\n", res.Name, msg.Topic())

	if res.DataType == "file" {
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
	} else {
		res, err := PrettyString(string(msg.Payload()))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
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

	var broker = "10.35.0.1"
	var port = 1883

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("mqtt-sub-cli")
	opts.SetUsername("cp_uqe_undef01")
	opts.SetPassword("FwWMuuVhx6Vgs8KV")
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
	topic := "cp/obd/cpa4000/#"
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic %s\n", topic)
}
