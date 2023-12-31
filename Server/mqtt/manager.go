package mqtt

import (
	"encoding/json"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Define constants for message types
const (
	Heartbeat               = 0
	HeartbeatAck            = 1
	SensorDiscovery         = 2
	SensorDiscoveryResponse = 3
	SensorData              = 4
)

// SensorInfo holds additional data about the sensor.
type SensorInfo struct {
	SensorName string `json:"sensor_name"` // Name/type of the sensor
	Units      string `json:"units"`       // Measurement units of the sensor
}

// DiscoveryMessage defines the structure of a sensor discovery message.
type DiscoveryMessage struct {
	Type      int        `json:"type"`      // Message type
	DeviceID  string     `json:"device_id"` // Unique identifier of the sensor
	Timestamp string     `json:"timestamp"` // Timestamp of the message
	Data      SensorInfo `json:"data"`      // Sensor information
}

type SensorStatus struct {
	SensorID   string `json:"sensor_id"`
	SensorName string `json:"sensor_name"`
	Units      string `json:"units"`
	Value      int    `json:"value"`
}

// ListenerManager manages a set of MQTT listeners for different sensors.
type ListenerManager struct {
	client         mqtt.Client          // MQTT client for communication
	listeners      map[string]*Listener // Map to hold listener instances indexed by device ID
	discoveryTopic string               // MQTT topic to listen to for sensor discovery
}

// A time constant used by the RemoveInactiveListeners method
const ListenerTimeout = time.Second * 30

// NewListenerManager creates and returns a new ListenerManager instance.
func NewListenerManager(client mqtt.Client, discoveryTopic string) *ListenerManager {
	return &ListenerManager{
		client:         client,
		listeners:      make(map[string]*Listener),
		discoveryTopic: discoveryTopic,
	}
}

// discoveryHandler handles incoming sensor discovery messages.
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
			// Create a new listener for the sensor if it doesn't already exist
			m.listeners[sensorID] = NewListener(m.client, sensorTopic, discoveryMsg.Data.SensorName, discoveryMsg.Data.Units)
			m.listeners[sensorID].Start()
			log.Printf("Listener started for sensor ID %s with name %s and units %s", sensorID, discoveryMsg.Data.SensorName, discoveryMsg.Data.Units)
		}
	}
}

// GetCurrentValues retrieves the current values from all listeners.
func (m *ListenerManager) GetCurrentValues() []SensorStatus {
	var sensorStatusList []SensorStatus

	for sensorId, listener := range m.listeners {
		sensorStatus := SensorStatus{
			SensorID:   sensorId,
			SensorName: listener.SensorName,
			Units:      listener.Units,
			Value:      listener.GetCurrentValue(),
		}
		sensorStatusList = append(sensorStatusList, sensorStatus)

	}
	return sensorStatusList
}

// Start begins the listening process on the discovery topic.
func (m *ListenerManager) Start() {
	if token := m.client.Subscribe(m.discoveryTopic, 0, m.discoveryHandler); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe to discovery topic %s: %v", m.discoveryTopic, token.Error())
	} else {
		log.Printf("Subscribed to discovery topic %s", m.discoveryTopic)
	}
}

func (m *ListenerManager) RemoveInactiveListeners() {
	for sensorID, listener := range m.listeners {
		if time.Since(listener.lastMessageTime) > ListenerTimeout {
			log.Printf("Removing inactive listener for sensor ID: %s", sensorID)
			delete(m.listeners, sensorID)
		}
	}
}
