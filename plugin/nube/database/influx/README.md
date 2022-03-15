# Scope
- Syncs data to InfluxDB2
- Flow: table.histories > InfluxDB2 
    - Organization: nube-org
    - Bucket: nube-bucket
- It uses config file for InfluxDB connection & job to sync that value
  - Use Token for credential; `influx auth list`
- It uses Job for starting scheduler task


### How to get default config
- Save empty `YAML` file, and it will generate you the default config file