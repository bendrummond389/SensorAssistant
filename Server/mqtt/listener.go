package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Listener represents an MQTT subscriber listening to a specific topic.
type Listener struct {
	client       mqtt.Client // MQTT client for communication
	topic        string      // MQTT topic the listener is subscribed to
	SensorType   string      // Type of sensor the listener is associated with
	Units        string      // Measurement units for sensor data
	currentValue int         // Most recent value received from the sensor
}

// SensorMessage defines the structure of a message from a sensor.
type SensorMessage struct {
	Type      int         `json:"type"`      // Message type
	DeviceID  string      `json:"device_id"` // Sensor's unique identifier
	Timestamp string      `json:"timestamp"` // Timestamp of the message
	Data      SensorValue `json:"data"`      // Sensor data
}

// SensorValue holds the actual data sent by a sensor.
type SensorValue struct {
	Value int `json:"value"` // Value reported by the sensor
}

// NewListener creates a new Listener instance.
func NewListener(client mqtt.Client, topic string, sensorType, units string) *Listener {
	return &Listener{
		client:       client,
		topic:        topic,
		SensorType:   sensorType,
		Units:        units,
		currentValue: 0,
	}
}

// GetCurrentValue returns the most recent value received from the sensor.
func (l *Listener) GetCurrentValue() int {
	return l.currentValue
}

// messageHandler handles incoming MQTT messages on the subscribed topic.
func (l *Listener) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))

	var sensorMsg SensorMessage
	if err := json.Unmarshal(msg.Payload(), &sensorMsg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	if sensorMsg.Type == SensorData {
		l.currentValue = sensorMsg.Data.Value
		log.Printf("Updated currentValue for topic %s: %d", l.topic, l.currentValue)
	}
}

// Start subscribes the Listener to its MQTT topic to start receiving messages.
func (l *Listener) Start() {
	if token := l.client.Subscribe(l.topic, 0, l.messageHandler); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe to topic %s, %v", l.topic, token.Error())
	} else {
		log.Printf("Subscribed to topic: %s", l.topic)
	}
}
