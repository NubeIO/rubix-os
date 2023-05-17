package shared

import (
	"context"
	"github.com/NubeIO/flow-framework/module/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

type DBHelper interface {
	Sum(int64, int64) (int64, error)
	CallAPI(path string, args string) ([]byte, error)
}

// Module is the interface that we're exposing as a plugin.
type Module interface {
	Init(dbHelper DBHelper) error
	Put(key string, value int64) error
	Get(key string) (int64, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
type NubeModule struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Module
}

func (p *NubeModule) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterModuleServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *NubeModule) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewModuleClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &NubeModule{}
