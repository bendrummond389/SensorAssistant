from mqtt.mqtt_manager import MQTTClientManager
import ntptime
import sensor
import time


if __name__ == "__main__":
    ntptime.settime()
    print("connecting client to broker")
    client = MQTTClientManager("Temperature Sensor", "Celsius")
    while True:
        data = sensor.read_sensor()
        client.publish_sensor_data(data)
        time.sleep(2)


