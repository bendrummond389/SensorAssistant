package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bendrummond389/SensorAssistant/Server/mqtt"
	"github.com/bendrummond389/SensorAssistant/Server/websocket"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	mqttBrokerAddress := os.Getenv("MQTT_BROKER_ADDRESS")
	if mqttBrokerAddress == "" {
		log.Printf("Broker address not provided from docker env, using local ip instead")
		mqttBrokerAddress = "tcp://localhost:1883"
	}
	opts := mqttPaho.NewClientOptions()
	opts.AddBroker(mqttBrokerAddress)

	client := mqttPaho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Error connecting to ")
	}
	manager := mqtt.NewListenerManager(client, "discovery")
	manager.Start()

	wsServer := websocket.NewServer()
	go wsServer.Run()

	http.HandleFunc("/ws", wsServer.HandleConnections)
	go http.ListenAndServe(":8080", nil)

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			sensorStatuses := manager.GetCurrentValues()
			data, err := json.Marshal(sensorStatuses)
			if err != nil {
				log.Println("Error marshaling sensor status", err)
				continue
			}

			wsServer.BroadcastToClients(data)
		}
	}()

	go func() {
		inactivityTicker := time.NewTicker(time.Second * 20)
		for range inactivityTicker.C {
			manager.RemoveInactiveListeners()
		}
	}()

	select {}
}
