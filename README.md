
# getting started
rename the `config-example.yml` file to `config.yml`

# default port
1660

# plugins
## See plugin docs
/docs/plugins


## Build plugin
add into /data/plugins

```
go build -buildmode=plugin -o ehco.so echo.go
go build -buildmode=plugin -o echo.so echo.go  && cp echo.so  ../../../data/plugins/  && go run ../../../app.go
```

example to build and run the apps
```
cd plugin/example/system
go build -buildmode=plugin -o system.so system.go  && cp system.so  ../../../data/plugins/ && (cd /home/aidan/code/go/flow-framework  && go run app.go  -config /home/aidan/code/go/flow-framework/config.yml)
```
