# Scope

- Syncs data to InfluxDB2
- Flow: table.histories > InfluxDB2 
    - Organization: nube-org
    - Bucket: nube-bucket
- It uses Integration for InfluxDB2 config & credentials
  - Use Token for credential; `influx auth list`
- It uses Job for starting scheduler task
