package runner

import (
	"context"
	"fmt"
	"github.com/NubeIO/rubix-os/config"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

func Run(engine http.Handler, conf *config.Configuration) {
	addr := fmt.Sprintf("%s:%d", conf.Server.ListenAddr, conf.Server.Port)
	log.Info("Started Listening for plain HTTP connection on " + addr)
	server := &http.Server{Addr: addr, Handler: engine}
	err := server.Serve(startListening(addr))
	log.Fatal(err)
}

func startListening(addr string) net.Listener {
	lc := net.ListenConfig{KeepAlive: -1}
	conn, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		log.Fatalln("Could not listen on", addr, err)
	}
	return conn
}
