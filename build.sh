#!/bin/bash

dir=$(pwd)
echo "$dir"
# shellcheck disable=SC2164
cd "$dir"/plugin/nube/system/
go build -buildmode=plugin -o system.so *.go  && cp system.so  "$dir"/data/plugins
# shellcheck disable=SC2164
cd "$dir"/plugin/nube/protocals/lorawan
go build -buildmode=plugin -o lorawan.so *.go  && cp lorawan.so "$dir"/data/plugins
# shellcheck disable=SC2164
cd "$dir"/plugin/nube/protocals/modbus
go build -buildmode=plugin -o modbus.so *.go  && cp modbus.so "$dir"/data/plugins
# shellcheck disable=SC2164
#cd "$dir"/plugin/nube/protocals/lora
#go build -buildmode=plugin -o lora.so *.go  && cp lora.so "$dir"/data/plugins
# shellcheck disable=SC2164
cd "$dir"/plugin/nube/protocals/bacnetserver
go build -buildmode=plugin -o bacnetserver.so *.go  && cp bacnetserver.so "$dir"/data/plugins
# shellcheck disable=SC2164
cd "$dir"
go run app.go


