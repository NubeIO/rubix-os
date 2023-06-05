package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/gin-gonic/gin"
)

type CloneEdgeDatabase interface {
	ResolveHost(uuid string, name string) (*model.Host, error)
	CloneEdge(globalUUID string, networks []*model.Network) error
}
type CloneEdgeApi struct {
	DB CloneEdgeDatabase
}

func (a *CloneEdgeApi) CloneEdge(ctx *gin.Context) {
	matchHostUUIDName(ctx)
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := client.NewClient(host.IP, host.Port, host.ExternalToken)
	networks, err := cli.GetNetworksForCloneEdge()
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	err = a.DB.CloneEdge(host.GlobalUUID, networks)
	ResponseHandler(nil, err, ctx)
}
