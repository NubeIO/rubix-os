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
    - Comparison Operators
        - `=`, `>`, `<`, `<=`, `>=`, `!=` 
    - Fields 
        - `value`
        - `timestamp`
        - `network_uuid`
        - `network_name`
        - `device_uuid`
        - `device_name`
        - `point_uuid`
        - `point_name`
        - `global_uuid`
        - `host_uuid`
        - `host_name`
        - `group_uuid`
        - `group_name`
        - `location_uuid`
        - `location_name`
        - `tag`
        - `meta_tag_key`
        - `meta_tag_value`
    - Filter examples:   
        ```
        1. network_name={network_name}&&device_name={device_name}&&point_name={point_name}&&host_uuid={host_uuid}
        2. (host_name={host_name}&&timestamp>{timestamp})||point_uuid!={point_uuid}
        3. (tag={tag}&&value>={value})||(tag={tag}&&value<={value})
        4. point_uuid=<{point_uuid}||point_uuid={point_uuid}
        5. (meta_tag_key={meta_tag_key}&&value>={value})
        6. (meta_tag_value={meta_tag_value}&&value<={value})
        7. (meta_tag_key={meta_tag_key}&&meta_tag_value={meta_tag_value}&&value>={value})
        ```
- limit
- offset
- order_by
- order [default: DESC]
- group_limit

### Endpoint
- GET `/api/plugins/api/postgres/histories?filter=<filter>&limit=<int>&offset=<int>&order_by=<order_by>&order=<order>&group_limit=<group_limit>`
