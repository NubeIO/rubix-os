package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

// https://www.chirpstack.io/

type ChirpProxyAPI struct {
}

func (inst *ChirpProxyAPI) ChirpProxy(c *gin.Context) { // eg http://0.0.0.0:8080/chrip/api/organizations?limit=10
	remote, err := Builder("0.0.0.0", 8080)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	token := c.GetHeader("cs-token")
	c.Request.Header.Del("host-uuid")
	c.Request.Header.Del("host-name")
	c.Request.Header.Del("authorization") // if this isn't deleted it cases issues on the cs server side
	c.Request.Header.Del("cs-token")
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
		if req.URL.Path != "/api/internal/login" {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Grpc-Metadata-Authorization", token) // pass in a header with the chirp-stack auth token
			log.Infof("chrip-proxy path:%s token-length: %d", req.URL.Path, len(token))
		}
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
