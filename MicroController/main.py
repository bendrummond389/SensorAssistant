# sensors/{sensor_type}/{sensor_id}/data for sending sensor data.
# sensors/{sensor_type}/{sensor_id}/status for sending sensor status updates.
# sensors/{sensor_type}/{sensor_id}/command for receiving commands.
# Consistent Naming Convention: Use a consistent and descriptive naming convention for topics to avoid confusion.

from mqtt.mqtt_manager import MQTTClientManager
import ntptime
import sensor
import time


if __name__ == "__main__":
    ntptime.settime()
    print("connecting client to broker")
    client = MQTTClientManager()
    while True:
        data = sensor.read_sensor()
        client.publish_sensor_data(data)
        time.sleep(5)
