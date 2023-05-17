package shared

import (
	"context"
	"github.com/NubeIO/flow-framework/module/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
)

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Module

	broker *plugin.GRPCBroker
}

func (m *GRPCServer) Init(ctx context.Context, req *proto.InitRequest) (*proto.Empty, error) {
	log.Info("gRPC Init server has been called...")
	conn, err := m.broker.Dial(req.AddServer)
	if err != nil {
		return nil, err
	}
	// defer conn.Close() // TODO: we haven't closed this
	dbHelper := &GRPCDBHelperClient{proto.NewDBHelperClient(conn)}
	return &proto.Empty{}, m.Impl.Init(dbHelper)
}

func (m *GRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	log.Info("gRPC put server has been called...")
	return &proto.Empty{}, m.Impl.Put(req.Key, req.Value)
}

func (m *GRPCServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := m.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCDBHelperClient struct{ client proto.DBHelperClient }

func (m *GRPCDBHelperClient) Sum(a, b int64) (int64, error) {
	resp, err := m.client.Sum(context.Background(), &proto.SumRequest{
		A: a,
		B: b,
	})
	if err != nil {
		hclog.Default().Info("add.Sum", "client", "start", "err", err)
		return 0, err
	}
	return resp.R, err
}

func (m *GRPCDBHelperClient) CallAPI(path, args string) ([]byte, error) {
	resp, err := m.client.CallAPI(context.Background(), &proto.APIRequest{
		Path: path,
		Args: args,
	})
	if err != nil {
		hclog.Default().Info("add.CallAPI", "client", "start", "err", err)
		return nil, err
	}
	return resp.R, err
}
