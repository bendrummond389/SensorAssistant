package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bendrummond389/SensorAssistant/Server/mqtt"
	"github.com/bendrummond389/SensorAssistant/Server/websocket"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	mqttBroker := "tcp://localhost:1883"
	opts := mqttPaho.NewClientOptions()
	opts.AddBroker(mqttBroker)

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

	ticker := time.NewTicker(10 * time.Second)
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

	select {}
}
