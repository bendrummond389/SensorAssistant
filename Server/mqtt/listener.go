package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Listener struct {
	client       mqtt.Client
	topic        string
	SensorType   string
	Units        string
	currentValue int
}

type SensorMessage struct {
	Type      int         `json:"type"`
	DeviceID  string      `json:"device_id"`
	Timestamp string      `json:"timestamp"`
	Data      SensorValue `json:"data"`
}

type SensorValue struct {
	Value int `json:"value"`
}

func NewListener(client mqtt.Client, topic string, sensorType, units string) *Listener {
	return &Listener{
		client:       client,
		topic:        topic,
		SensorType:   sensorType,
		Units:        units,
		currentValue: 0,
	}
}

func (l *Listener) GetCurrentValue() int {
	return l.currentValue
}

func (l *Listener) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Recieved message on topic %s: %v", msg.Topic(), string(msg.Payload()))

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

func (l *Listener) Start() {
	if token := l.client.Subscribe(l.topic, 0, l.messageHandler); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe to topic %s, %v", l.topic, token.Error())
	} else {
		log.Printf("Subscribed to topic: %s", l.topic)
	}
}
