import dht
import machine

def read_sensor():
    try:
        d = dht.DHT11(machine.Pin(2))
        d.measure()
        val = int(d.temperature())
        print(f"Sensor read successfully: {val}")
        return val
    except Exception as e:
        print(f"Exception in read_sensor: {e}")
        return None
