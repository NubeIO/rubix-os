# getting started

rename the `config-example.yml` file to `config.yml`

run the bash script to build and start with plugins
`bash build.bash --system --rubix --modbus --lora`

# default port

1660

# plugins

## See plugin docs

/docs/plugins

## Build plugin

add into /data/plugins

```
go build -buildmode=plugin -o ehco.so *.go
```

example to build and run the apps

```
cd plugin/example/system
go build -buildmode=plugin -o system.so *.go  && cp system.so  ../../../data/plugins/ && rm system.so && (cd ~/code/go/nube/flow-framework  && go run app.go)
```

## Logging

```
debug: when we want to show information on debugging issue (we activate this mode on just debugging so will not be that much un-necessary logs)
info: when we want to show meaningful information for user
warn: when we want to give a warning for user for some operations
error: while error happens, show it on red alert  
```

## APIs endpoints

```
- consumers with args ?writers=true
- producers with args ?writers=true
- flow_networks with args ?streams=true
- child params args available for networks, devices ?devices=true&points=true&serial_connection=true&ip_connection=true

- /api/flow/networks?global_uuid=global_uuid1
- /api/flow/networks/one/args?global_uuid=global_uuid1&client_id=client_id1&site_id=site_id1&device_id=device_id1
```
