# Scope

- Syncs data to PostgreSQL
- Flow: table.histories > PostgreSQL
- It uses config file for PostgreSQL connection & job to sync that value
- It uses Job for starting scheduler task

### How to get default config

- Save empty `YAML` file, and it will generate you the default config file

# History api

### Query params
- filter
    - Logical Operators `(make sure it is url encoded)`
        - `&&`, `||` 
    - Comparision Operators
        - `=`, `>`, `<`, `<=`, `>=`, `!=` 
    - Fields 
        - `value`
        - `timestamp`
        - `rubix_network_uuid`
        - `rubix_network_name`
        - `rubix_device_uuid`
        - `rubix_device_name`
        - `rubix_point_uuid`
        - `rubix_point_name`
        - `global_uuid`
        - `client_id`
        - `client_name`
        - `site_id`
        - `site_name`
        - `device_id`
        - `device_name`
        - `tag`     
    - Filter examples:   
        ```
        1. rubix_network_name={network_name}&&rubix_device_name={device_name}&&rubix_point_name={point_name}&&site_id={site_id}
        2. (site_name={site_name}&&timestamp>{timestamp})||rubix_point_uuid!={point_uuid}
        3. (tag={tag}&&value>={value})||(tag={tag}&&value<={value})
        4. rubix_point_uuid=<{oint_uuid}||rubix_point_uuid={point_uuid}
        ```
- limit
- offset

### Endpoint
- GET `/api/plugins/api/postgres/histories?filter=<filter>&limit=<int>&offset=<int>`
