# scope

- Creates chirpstack devices (application must be created already, and application number into plugin config)
- Creates map.txt file for LoRaWAN <=> Modbus Bridge
- Sets the names of TMV points (from lorawan defaults).
- Creates Modbus Network, Devices, and Points from auto added lorawan devices and points
- Runs RTC Sync Daily

Sample Config for Commissioning:
```
job:
  enable_config_steps: true
  frequency: 30m
  chirpstack_application_number: 1
  chirpstack_network_key: 0301021604050F07E6095A0B0C12630F
  chirpstack_username: admin
  chirpstack_password: Helensburgh2508
  tmv_json_file_path: /home/pi/test.json
log_level: DEBUG
```

Sample Config for Ongoing Run (RTC Sync Only):
```
job:
  enable_config_steps: false
  frequency: 30m
  chirpstack_application_number: 1
  chirpstack_network_key: 0301021604050F07E6095A0B0C12630F
  chirpstack_username: admin
  chirpstack_password: Helensburgh2508
  tmv_json_file_path: /home/pi/test.json
log_level: ERROR
```