package api

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/ip"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
)

func composeExternalToken(token string) string {
	return fmt.Sprintf("External %s", token)
}

type HostProxyDatabase interface {
	ResolveHost(uuid string, name string) (*model.Host, error)
}

type HostProxyAPI struct {
	DB HostProxyDatabase
}

func (a *HostProxyAPI) HostProxy(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	proxyPath := strings.Trim(ctx.Param("proxyPath"), string(os.PathSeparator))
	proxyPathParts := strings.Split(proxyPath, "/")
	var remote *url.URL = nil
	if len(proxyPathParts) > 0 && proxyPathParts[0] == "eb" {
		proxyPath = path.Join(proxyPathParts[1:]...)
		remote, err = ip.Builder(host.HTTPS, host.IP, host.BiosPort)
	} else if len(proxyPathParts) > 0 && proxyPathParts[0] == "ros" {
		proxyPath = path.Join(proxyPathParts[1:]...)
		remote, err = ip.Builder(host.HTTPS, host.IP, host.Port)
	} else {
		remote, err = ip.Builder(host.HTTPS, host.IP, host.Port)
	}
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	externalToken := host.ExternalToken
	proxyPath = fmt.Sprintf("/%s", proxyPath)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = proxyPath
		authorization := ctx.GetHeader("jwt-token")
		if authorization != "" {
			req.Header.Set("Authorization", authorization)
		} else if externalToken != "" {
			req.Header.Set("Authorization", composeExternalToken(externalToken))
		}
	}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
