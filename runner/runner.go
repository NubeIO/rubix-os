package runner

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/NubeIO/flow-framework/config"
)

func Run(engine *gin.Engine, conf *config.Configuration) {
	addr := fmt.Sprintf("%s:%d", conf.Server.ListenAddr, conf.Server.Port)
	log.Info("Started Listening for plain HTTP connection on " + addr)
	pprof.Register(engine)
	adminGroup := engine.Group("/admin", func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != "foobar" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	})
	pprof.RouteRegister(adminGroup, "pprof")
	engine.Run(addr)
}
