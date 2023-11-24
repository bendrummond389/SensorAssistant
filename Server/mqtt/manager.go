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

type SensorInfo struct {
	SensorName string `json:"sensor_name"`
	Units      string `json:"units"`
}

type DiscoveryMessage struct {
	Type      int        `json:"type"`
	DeviceID  string     `json:"device_id"`
	Timestamp string     `json:"timestamp"`
	Data      SensorInfo `json:"data"`
}

type ListenerManager struct {
	client         mqtt.Client
	listeners      map[string]*Listener
	discoveryTopic string
}

func NewListenerManager(client mqtt.Client, discoveryTopic string) *ListenerManager {
	return &ListenerManager{
		client:         client,
		listeners:      make(map[string]*Listener),
		discoveryTopic: discoveryTopic,
	}
}

func (m *ListenerManager) discoveryHandler(client mqtt.Client, msg mqtt.Message) {
	var discoveryMsg DiscoveryMessage
	if err := json.Unmarshal(msg.Payload(), &discoveryMsg); err != nil {
		log.Printf("Error unmarshaling discovery message: %v", err)
		return
	}

	if discoveryMsg.Type == SensorDiscovery {
		sensorID := discoveryMsg.DeviceID
		sensorTopic := sensorID + "/data"
		if _, exists := m.listeners[sensorID]; !exists {
			m.listeners[sensorID] = NewListener(m.client, sensorTopic, discoveryMsg.Data.SensorName, discoveryMsg.Data.Units)
			m.listeners[sensorID].Start()
			log.Printf("Listener started for sensor ID %s with name %s and units %s", sensorID, discoveryMsg.Data.SensorName, discoveryMsg.Data.Units)
		}
	}
}

func (m *ListenerManager) GetCurrentValues() map[string]int {
	currentValues := make(map[string]int)
	for sensorID, listener := range m.listeners {
		currentValues[sensorID] = listener.GetCurrentValue()
	}
	return currentValues
}

func (m *ListenerManager) Start() {
	if token := m.client.Subscribe(m.discoveryTopic, 0, m.discoveryHandler); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe to discovery topic %s: %v", m.discoveryTopic, token.Error())
	} else {
		log.Printf("Subscribed to discovery topic %s", m.discoveryTopic)
	}
}
