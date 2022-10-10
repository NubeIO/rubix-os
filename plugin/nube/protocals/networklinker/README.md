# Combines points from LoRaRAW/WAN and Modbus

Used for devices that push data over LoRaRAW/WAN but receive writes via modbus

### Scope

- All API writes are directed towards modbus
- LoRa points do not write to modbus to avoid unnecessary writes and looping