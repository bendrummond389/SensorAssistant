# Import necessary libraries
import network
import time
import json


# Function to load MQTT configuration from a file
def load_connection_settings():
    try:
        # Attempt to open and read the configuration file
        with open("connection_settings.json", "r") as file:
            return json.load(file)
    except Exception as e:
        # Print an error message if the file could not be read
        print(f"Could not read MQTT config: {e}")
        return {}  # Return an empty dictionary on failure


# Load the MQTT configuration
config = load_connection_settings()

# Get the SSID and WiFi password from the configuration, or use default values
SSID = config.get("SSID", "default")
WIFI_PASSWORD = config.get("WIFI_PASSWORD", "default")


# Function to connect to the WiFi network
def connect_wifi():
    # Create a WLAN station interface
    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)  # Activate the interface
    ssid = SSID
    password = WIFI_PASSWORD

    retries = 5  # Set the number of retry attempts for connection

    # Loop until connected or retries are exhausted
    while not wlan.isconnected() and retries > 0:
        print("Attempting to connect to network...")

        if not wlan.isconnected():
            wlan.connect(ssid, password)  # Try to connect to the network

            # Check the connection status for 10 seconds
            for i in range(10):
                if wlan.isconnected():
                    break  # Break the loop if connected
                time.sleep(1)  # Sleep for a second between checks

        retries -= 5  # Decrement the retry counter

        # Print the connection status
        if wlan.isconnected():
            print(f"connected to {ssid}")
            print("Network config:", wlan.ifconfig())  # Print network configuration
            break  # Break the loop if connected
        else:
            print(
                f"Failed to connect. Retries left: {retries}"
            )  # Print the number of retries left if connection failed


# Call the function to connect to WiFi
connect_wifi()
