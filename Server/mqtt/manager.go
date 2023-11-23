package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	Heartbeat               = 0
	HeartbeatAck            = 1
	SensorDiscovery         = 2
	SensorDiscoveryResponse = 3
	SensorData              = 4
)

type Message struct {
	Type      int    `json:"type"`
	DeviceID  string `json:"device_id"`
	Timestamp string `json:"timestamp"`
	Data      interface{} `json:"data"`
}


type ListenerManager struct {
	client mqtt.Client
	listeners map[string]*Listener
	discoveryTopic string
}


func NewListenerManager(client mqtt.Client, discoveryTopic string) *ListenerManager {
	return &ListenerManager{
		client: client,
		listeners: make(map[string]*Listener),
		discoveryTopic: discoveryTopic,

	}
}

func (m *ListenerManager) discoveryHandler(client mqtt.Client, msg mqtt.Message) {
	var message Message
	if err := json.Unmarshal(msg.Payload(), &message); err != nil {
		log.Printf("Error unmarshaling discovery message: %v", err)
		return
	}

	if message.Type == SensorDiscovery {
		sensorID := message.DeviceID
		sensorTopic := sensorID + "/data"
		if _, exists := m.listeners[sensorID]; !exists {
			m.listeners[sensorID] = NewListener(m.client, sensorTopic)
			m.listeners[sensorID].Start()
			log.Printf("Listener started for topic: %s", sensorID)
		}
	}


}

func (m *ListenerManager) Start() {
	if token := m.client.Subscribe(m.discoveryTopic, 0, m.discoveryHandler); token.Wait() && token.Error() != nil {
			log.Printf("Failed to subscribe to discovery topic %s: %v", m.discoveryTopic, token.Error())
	} else {
			log.Printf("Subscribed to discovery topic %s", m.discoveryTopic)
	}
}
