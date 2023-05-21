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
	GetWithoutParam(path, args string) ([]byte, error)
	Get(path, uuid, args string) ([]byte, error)
	Post(path string, body []byte) ([]byte, error)
	Put(path, uuid string, body []byte) ([]byte, error)
	Patch(path, uuid string, body []byte) ([]byte, error)
	Delete(path, uuid string) ([]byte, error)
}

type Info struct {
	Name       string
	Author     string
	Website    string
	License    string
	HasNetwork bool
}

// Module is the interface that we're exposing as a plugin.
type Module interface {
	Init(dbHelper DBHelper, moduleName string) error
	GetInfo() (*Info, error)
	GetUrlPrefix() (*string, error)
	Get(path string) ([]byte, error)
	Post(path string, body []byte) ([]byte, error)
	Put(path string, body []byte) ([]byte, error)
	Patch(path string, body []byte) ([]byte, error)
	Delete(path string) ([]byte, error)
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