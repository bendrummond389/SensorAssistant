# Sensor Assistant

## Overview

This project provides a framework to create a personal network of sensors, coupled with a simple web-based interface for viewing sensor readings. It's designed to lower the barrier to IoT exploration, enabling users to set up a home sensor network with just a basic understanding of Python.

## Features

- **Custom Sensor Network**: Build your own network with various sensors.
- **Web Interface**: Easily view sensor data through a user-friendly web interface.
- **User-Friendly**: Requires only basic Python knowledge for setup and customization.

## Getting Started

### Prerequisites

- Python 3.x
- Basic knowledge of Python programming
- Required hardware sensors (list specific models or types)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/sensornetworkframework.git
   ```

2. Navigate to the project directory:

   ```bash
   cd sensornetworkframework
   ```

### Configuration

- Modify three key files in this project.
- First, add the SSID and password of your network to `connection_settings.py`.
- Next, add the IP address of the server device to `mqtt/mqtt_config`.
- Finally, modify the `read_sensor()` function in `sensor.py` to return the desired value from your sensor. Note: There's currently an example for reading a water sensor. More examples and possibly wiring diagrams may be added in the future.

### Running the Application

The application currently requires a few steps to run. To simplify the process the MQTT broker, go server and nextjs frontend have been containerized and can be built and deployed with 2 commands.

1. **Deploy MicroPython to Your Microcontroller**:

   - Use pymaker for easy file flashing to the device. Once deployed, the code runs on startup, so no further configuration is needed.

2. **Build the Docker image**:

   - From the project root, run:

     ```bash
      docker-compose build
     ```

3. **Start the docker image**:

   - From the project root, run:

     ```bash
      docker compose up
     ```

## Usage

Once the application is started, visit [localhost:3000](http://localhost:3000) to view the web-based dashboard.

## License

This project is licensed under the [MIT License](LICENSE) - see the LICENSE file for details.
