# sensors/{sensor_type}/{sensor_id}/data for sending sensor data.
# sensors/{sensor_type}/{sensor_id}/status for sending sensor status updates.
# sensors/{sensor_type}/{sensor_id}/command for receiving commands.
# Consistent Naming Convention: Use a consistent and descriptive naming convention for topics to avoid confusion.

from mqtt.mqtt_manager import MQTTClientManager




if __name__ == "__main__":
    print("we made it to main")
    
    client = MQTTClientManager()
    client.connect_to_broker()
    
