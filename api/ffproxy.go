package api

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/internaltoken"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Builder(ip string, port int) (*url.URL, error) {
	return url.ParseRequestURI(CheckHTTP(fmt.Sprintf("%s:%d", ip, port)))
}

func CheckHTTP(address string) string {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		return "http://" + address
	}
	return address
}

type FFProxyAPI struct {
}

func (inst *FFProxyAPI) FFProxy(c *gin.Context) {
	remote, err := Builder("0.0.0.0", 1660)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	internalToken := internaltoken.GetInternalToken(true)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
		req.Header.Set("Authorization", internalToken)
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
