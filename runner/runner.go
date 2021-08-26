package runner

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/NubeDev/flow-framework/config"
)

func Run(router http.Handler, conf *config.Configuration) {
	httpHandler := router
	addr := fmt.Sprintf("%s:%d", conf.Server.ListenAddr, conf.Server.Port)
	log.Println("Started Listening for plain HTTP connection on " + addr)
	server := &http.Server{Addr: addr, Handler: httpHandler}
	log.Fatal(server.Serve(startListening(addr, conf.Server.KeepAlivePeriodSeconds)))
}

func startListening(addr string, keepAlive int) net.Listener {
	lc := net.ListenConfig{KeepAlive: time.Duration(keepAlive) * time.Second}
	conn, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		log.Fatalln("Could not listen on", addr, err)
	}
	return conn
}
