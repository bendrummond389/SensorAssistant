package main

import (
	"log"
	"time"
	"github.com/bendrummond389/SensorAssistant/Server/mqtt"
	mqttPaho "github.com/eclipse/paho.mqtt.golang"
	// other imports
)

func main() {
	mqttBroker := "tcp://localhost:1883"
	opts := mqttPaho.NewClientOptions()
	opts.AddBroker(mqttBroker)

	client := mqttPaho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("Error connecting to ")
	}

	topic := "discovery"
	manager := mqtt.NewListenerManager(client, topic)
	manager.Start()

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			currentValues := manager.GetCurrentValues()
			log.Println("Current Sensor Values:", currentValues)
		}
	}()

	select {}
}