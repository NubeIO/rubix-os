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
  <client_id>/<client_name>/<site_id>/<site_name>/<device_id>/<device_name>/rubix/platform/points/publish
```

**get device platform info:**

send a message to these topic to get the device info

get the device platform info (will return all the platform info in rubix-os)

```
rubix/platform/info
// will response on this topic
rubix/platform/info/publish
```

get the edge-device points list (will return all the points in rubix-os)

```
rubix/platform/points
// will response on this topic
rubix/platform/points/publish
```

get the edge-device point (will return a single point in rubix-os)

```
rubix/platform/point
body by name
{
    "network_name": "net",
    "device_name": "dev",
    "point_name": "pnt"
}
body by uuid
{
    "point_uuid": "pnt_94ea3ea254dc440a"
}
// will response on this topic
rubix/platform/point/publish
```

point write over MQTT

```
rubix/platform/points/write
body by name 
{
    "network_name": "net",
    "device_name": "dev",
    "point_name": "pnt",
    "priority": {
        "priority": {
            "_1": null,
            "_2": null,
            "_3": null,
            "_4": null,
            "_5": null,
            "_6": null,
            "_7": null,
            "_8": null,
            "_9": null,
            "_10": null,
            "_11": null,
            "_12": null,
            "_13": null,
            "_14": null,
            "_15": null,
            "_16": 11345
        }
    }
}
body by uuid
{
    "point_uuid": "pnt_94ea3ea254dc440a"
     "priority": {
        "priority": {
            "_1": null,
            "_2": null,
            "_3": null,
            "_4": null,
            "_5": null,
            "_6": null,
            "_7": null,
            "_8": null,
            "_9": null,
            "_10": null,
            "_11": null,
            "_12": null,
            "_13": null,
            "_14": null,
            "_15": null,
            "_16": 11345
        }
    }
}
// will response on this topic
rubix/platform/list/points/publish
```


