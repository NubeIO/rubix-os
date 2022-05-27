# Scope

- Syncs data to PostgreSQL
- Flow: table.histories > PostgreSQL
- It uses config file for PostgreSQL connection & job to sync that value
- It uses Job for starting scheduler task

### How to get default config

- Save empty `YAML` file, and it will generate you the default config file