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

// Here is the gRPC server that GRPCClient talks to.
type GRPCDBHelperServer struct {
	// This is the real implementation
	Impl DBHelper
}

func (m *GRPCDBHelperServer) GetList(ctx context.Context, req *proto.GetListRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.GetList(req.Path, req.Args)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCDBHelperServer) Get(ctx context.Context, req *proto.GetRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Get(req.Path, req.Uuid, req.Args)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCDBHelperServer) Post(ctx context.Context, req *proto.PostRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Post(req.Path, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCDBHelperServer) Put(ctx context.Context, req *proto.PutRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Put(req.Path, req.Uuid, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCDBHelperServer) Patch(ctx context.Context, req *proto.PatchRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Patch(req.Path, req.Uuid, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCDBHelperServer) Delete(ctx context.Context, req *proto.DeleteRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Delete(req.Path, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}
