import json
import time
from simple import MQTTClient
import machine
import ubinascii


class MQTTClientManager:
    def __init__(self, config_file="mqtt/mqtt_config.json"):
        self.client = None
        self.device_id = ubinascii.hexlify(machine.unique_id())
        self.broker_address = ""
        self.broker_port = 0
        self.load_mqtt_config(config_file)

    def load_mqtt_config(self, config_file):
        try:
            with open(config_file, "r") as file:
                config = json.load(file)
                self.broker_address = config.get("BROKER_ADDRESS", "default_broker")
                self.broker_port = config.get("BROKER_PORT", 1883)
        except Exception as e:
            print(f"Could not read MQTT config: {e}")

    def connect_to_broker(self):
        try:
            self.client = MQTTClient(
                self.device_id, self.broker_address, port=self.broker_port
            )
            self.client.set_callback(self.mqtt_callback)
            self.client.connect
            print(f"connected to broker as {self.device_id}")

        except Exception as e:
            print(f"Exception during MQTT connection: {e}")
            
    def send_device_info_to_discovery(self):
        self.client.publish()

    def mqtt_callback(self, topic, msg):
        print(f"print message recieved on topic {topic}: {msg}")

    # def publish_data(self, data):
    #     topic = f"sensors/{self.device_id}/data"
    #     payload = json.dumps({

    #     })
