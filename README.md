# getting started

rename the `config-eg.yml` file to `config.yml`

run the bash script to build and start with plugins
```
bash build.bash --help
```

## Default Port

1660

## Plugins

### See plugin docs

/docs/plugins

## Logging

```
debug: when we want to show information on debugging issue (we activate this mode on just debugging so will not be that much un-necessary logs)
info: when we want to show meaningful information for user
warn: when we want to give a warning for user for some operations
error: while error happens, show it on red alert  
```

### MQTT client
  
#### Topic structure:
```
<client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/points/value/<event>/...
<client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/points/value/cov/all/<network_plugin_path>/<network_uuid>/<network_name>/<device_uuid>/<device_name>/<point_uuid>/<point_name>
```
```
COV:
<client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/points/value/cov/all/<network_plugin_path>/<network_uuid>/<network_name>/<device_uuid>/<device_name>/<point_uuid>/<point_name>
```


#### Example topics:

**COV:**
```
all points:
  <client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/points/value/cov/all/#

by network plugin path:
  <client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/points/value/cov/all/<network_plugin_path>/+/+/+/+/+/+

by point uuid:
  <client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/points/value/cov/all/+/+/+/+/+/<point_uuid>/+

by point name:
  <client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/points/value/cov/all/+/+/<network_name>/+/<device_name>/+/<point_name>
```
**List:**
```
points list:
  <client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/points/value/points
```

**get device platform info:**

send a message to these topic to get the device info

get the device platform info
```
rubix/platform/info
// will response on this topic
rubix/platform/info/publish
```

get the edge-device points list
```
rubix/platform/points
// will response on this topic
rubix/platform/points/publish
```


  