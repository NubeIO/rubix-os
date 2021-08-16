



Build plugin
add into /data/plugins

```
go build -buildmode=plugin -o ehco.so echo.go
go build -buildmode=plugin -o echo.so echo.go  && cp echo.so  ../../../data/plugins/ 

```
