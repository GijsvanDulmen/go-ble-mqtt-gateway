# Golang BLE to MQTT Gateway

Simple golang based BLE to MQTT gateway. Sends all information it can retrieve with just listening for Bluetooth Low Energy advertisements. Works on OSX and Linux.

# Usage example

```
export MQTT_BROKER="tcp://192.168.1.244:1883"
export MQTT_CLIENT_ID="ble2mqtt-gw"
export MQTT_USERNAME="insert username"
export MQTT_PASSWORD="insert broker password"
export BLE_ID="BLEGW01"
export BLE_TOPIC="ble2mqtt/"
export BLE_NAME="BLE Scanner in room X"
export LOGGING="true"
go run .
```