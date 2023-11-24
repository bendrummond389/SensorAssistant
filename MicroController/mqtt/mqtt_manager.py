import json
import time
from simple import MQTTClient
import machine
import ubinascii

(
    HEARTBEAT,
    HEARTBEAT_ACK,
    SENSOR_DISCOVERY,
    SENSOR_DISCOVERY_RESPONSE,
    SENSOR_DATA,
) = range(5)


class MQTTClientManager:
    def __init__(self, name, units, config_file="mqtt/mqtt_config.json"):
        self.client = None
        self.name = name
        self.units = units
        self.device_id = ubinascii.hexlify(machine.unique_id()).decode("utf-8")
        self.broker_address = ""
        self.broker_port = 0
        self.load_mqtt_config(config_file)
        self.connect_to_broker()
        self.send_device_info_to_discovery()

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
            self.client.connect()
            print(f"connected to broker as {self.device_id}")

        except Exception as e:
            print(f"Exception during MQTT connection: {e}")

    def send_device_info_to_discovery(self):
        try:
            sensor_info = {
                "sensor_name": self.name,
                "units": self.units,
            }
            payload = self.create_payload(SENSOR_DISCOVERY, data=sensor_info)
            self.client.publish(topic="discovery", msg=payload)
        except Exception as e:
            print(f"Exception while sending device info to the discovery: {e}")

    def mqtt_callback(self, topic, msg):
        print(f"print message recieved on topic {topic}: {msg}")

    def create_payload(self, msg_type, data=None):
        return json.dumps(
            {
                "device_id": self.device_id,
                "type": msg_type,
                "timestamp": self.iso_timestamp(),
                "data": data,
            }
        )

    def publish_sensor_data(self, value):
        sensor_data = {"value": value}
        payload = self.create_payload(msg_type=SENSOR_DATA, data=sensor_data)
        self.client.publish(f"{self.device_id}/data", payload)

    def iso_timestamp(self):
        year, month, day, hour, minute, second, _, _ = time.localtime()
        return "{:04d}-{:02d}-{:02d}T{:02d}:{:02d}:{:02d}Z".format(
            year, month, day, hour, minute, second
        )
