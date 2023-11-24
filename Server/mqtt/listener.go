package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Listener struct {
	client mqtt.Client
	topic string
	SensorType string
	Units string
}

func NewListener(client mqtt.Client, topic string, sensorType, units string) *Listener {
	return &Listener{
		client: client,
		topic: topic,
		SensorType: sensorType,
		Units: units,
	}
}

func (l *Listener) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Recieved message on topic %s: %v", msg.Topic(), string(msg.Payload()))

	var data interface{}
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
	}

}

func (l *Listener) Start() {
	if token := l.client.Subscribe(l.topic, 0, l.messageHandler); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe to topic %s, %v", l.topic, token.Error())
	} else {
		log.Printf("Subscribed to topic: %s", l.topic)
	}
}