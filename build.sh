#!/bin/bash

pwd=`pwd`

cd "$pwd"/plugin/nube/system
go build -buildmode=plugin -o system.so *.go  && cp system.so  ../../../data/plugins/ && rm system.so
cd /home/aidan/code/go/nube/flow-framework/plugin/nube/protocals/lorawan
go build -buildmode=plugin -o lorawan.so *.go  && cp lorawan.so  ../../../../data/plugins/ && rm lorawan.so && (cd ~/code/go/nube/flow-framework  && go run app.go)
#cd /home/aidan/code/go/nube/flow-framework/plugin/nube/protocals/bacnetserver
#go build -buildmode=plugin -o bacnet.so *.go  && cp bacnet.so  ../../../../data/plugins/ && rm bacnet.so && (cd ~/code/go/nube/flow-framework  && go run app.go)

