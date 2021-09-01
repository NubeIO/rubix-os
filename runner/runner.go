package runner

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"

	"github.com/NubeDev/flow-framework/config"
)

func Run(engine http.Handler, conf *config.Configuration) {
	addr := fmt.Sprintf("%s:%d", conf.Server.ListenAddr, conf.Server.Port)
	log.Info("Started Listening for plain HTTP connection on " + addr)
	server := &http.Server{Addr: addr, Handler: engine}
	err := server.Serve(startListening(addr, conf.Server.KeepAlivePeriodSeconds))
	log.Fatal(err)
}

func startListening(addr string, keepAlive int) net.Listener {
	lc := net.ListenConfig{KeepAlive: time.Duration(keepAlive) * time.Second}
	conn, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		log.Fatalln("Could not listen on", addr, err)
	}
	return conn
}
