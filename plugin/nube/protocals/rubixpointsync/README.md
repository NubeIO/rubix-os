# Scope

- Syncs points from Flow Framework to Rubix Point Server Points

### How to get default config

- Save empty `YAML` file, and it will generate you the default config file

### Influx v1 config

```
influx:
- host: <host with no http/https>
  port: <port>
  token: <database name>:<password>
  org: ""  <-- leave this blank
  bucket: <database name>
  measurement: <add a name here>
job:
  frequency: 5m
  networks:
  - system
log_level: DEBUG
```