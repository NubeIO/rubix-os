package shared

import (
	"context"
	"github.com/NubeIO/flow-framework/module/proto"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCClient is an implementation of Module that talks over RPC.
type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.ModuleClient
}

func (m *GRPCClient) Init(dbHelper DBHelper) error {
	log.Infof("gRPC Init client has been called...")
	dbHelperServer := &GRPCDBHelperServer{Impl: dbHelper}
	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterDBHelperServer(s, dbHelperServer)

		return s
	}
	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	_, err := m.client.Init(context.Background(), &proto.InitRequest{
		AddServer: brokerID,
	})

	// s.Stop() // TODO: we haven't closed this
	return err
}

func (m *GRPCClient) Put(key string, value int64) error {
	log.Info("gRPC put client has been called...")
	_, err := m.client.Put(context.Background(), &proto.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (m *GRPCClient) Get(key string) (int64, error) {
	resp, err := m.client.Get(context.Background(), &proto.GetRequest{
		Key: key,
	})
	if err != nil {
		return 0, err
	}

	return resp.Value, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCDBHelperServer struct {
	// This is the real implementation
	Impl DBHelper
}

func (m *GRPCDBHelperServer) Sum(ctx context.Context, req *proto.SumRequest) (resp *proto.SumResponse, err error) {
	r, err := m.Impl.Sum(req.A, req.B)
	if err != nil {
		return nil, err
	}
	return &proto.SumResponse{R: r}, err
}

func (m *GRPCDBHelperServer) CallAPI(ctx context.Context, req *proto.APIRequest) (resp *proto.APIResponse, err error) {
	r, err := m.Impl.CallAPI(req.Path, req.Args)
	if err != nil {
		return nil, err
	}
	return &proto.APIResponse{R: r}, err
}
