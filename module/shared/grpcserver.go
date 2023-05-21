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
	log.Debug("gRPC Init server has been called...")
	conn, err := m.broker.Dial(req.AddServer)
	if err != nil {
		return nil, err
	}
	// defer conn.Close() // TODO: we haven't closed this
	dbHelper := &GRPCDBHelperClient{proto.NewDBHelperClient(conn)}
	return &proto.Empty{}, m.Impl.Init(dbHelper, req.ModuleName)
}

func (m *GRPCServer) GetInfo(ctx context.Context, req *proto.Empty) (*proto.InfoResponse, error) {
	log.Debug("gRPC GetInfo server has been called...")
	r, err := m.Impl.GetInfo()
	if err != nil {
		return nil, err
	}
	return &proto.InfoResponse{
		Name:       r.Name,
		Author:     r.Author,
		Website:    r.Website,
		License:    r.License,
		HasNetwork: r.HasNetwork,
	}, err
}

func (m *GRPCServer) GetUrlPrefix(ctx context.Context, req *proto.Empty) (*proto.UrlPrefixResponse, error) {
	log.Debug("gRPC GetUrlPrefix server has been called...")
	r, err := m.Impl.GetUrlPrefix()
	if err != nil {
		return nil, err
	}
	return &proto.UrlPrefixResponse{R: *r}, err
}

func (m *GRPCServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.Response, error) {
	log.Debug("gRPC Get server has been called...")
	r, err := m.Impl.Get(req.Path)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCServer) Post(ctx context.Context, req *proto.PostRequest) (*proto.Response, error) {
	log.Debug("gRPC Post server has been called...")
	r, err := m.Impl.Post(req.Path, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Response, error) {
	log.Debug("gRPC Put server has been called...")
	r, err := m.Impl.Put(req.Path, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCServer) Patch(ctx context.Context, req *proto.PatchRequest) (*proto.Response, error) {
	log.Debug("gRPC Patch server has been called...")
	r, err := m.Impl.Patch(req.Path, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

func (m *GRPCServer) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.Response, error) {
	log.Debug("gRPC Delete server has been called...")
	r, err := m.Impl.Delete(req.Path)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, err
}

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCDBHelperClient struct{ client proto.DBHelperClient }

func (m *GRPCDBHelperClient) GetWithoutParam(path, args string) ([]byte, error) {
	resp, err := m.client.GetWithoutParam(context.Background(), &proto.GetWithoutParamRequest{
		Path: path,
		Args: args,
	})
	if err != nil {
		hclog.Default().Info("GetList", err)
		return nil, err
	}
	return resp.R, err
}

func (m *GRPCDBHelperClient) Get(path, uuid, args string) ([]byte, error) {
	resp, err := m.client.Get(context.Background(), &proto.GetRequest{
		Path: path,
		Uuid: uuid,
		Args: args,
	})
	if err != nil {
		hclog.Default().Info("Get", err)
		return nil, err
	}
	return resp.R, err
}

func (m *GRPCDBHelperClient) Post(path string, body []byte) ([]byte, error) {
	resp, err := m.client.Post(context.Background(), &proto.PostRequest{
		Path: path,
		Body: body,
	})
	if err != nil {
		hclog.Default().Info("Post", err)
		return nil, err
	}
	return resp.R, err
}

func (m *GRPCDBHelperClient) Put(path, uuid string, body []byte) ([]byte, error) {
	resp, err := m.client.Put(context.Background(), &proto.PutRequest{
		Path: path,
		Uuid: uuid,
		Body: body,
	})
	if err != nil {
		hclog.Default().Info("Put", err)
		return nil, err
	}
	return resp.R, err
}

func (m *GRPCDBHelperClient) Patch(path, uuid string, body []byte) ([]byte, error) {
	resp, err := m.client.Patch(context.Background(), &proto.PatchRequest{
		Path: path,
		Uuid: uuid,
		Body: body,
	})
	if err != nil {
		hclog.Default().Info("Patch", err)
		return nil, err
	}
	return resp.R, err
}

func (m *GRPCDBHelperClient) Delete(path, uuid string) ([]byte, error) {
	resp, err := m.client.Delete(context.Background(), &proto.DeleteRequest{
		Path: path,
		Uuid: uuid,
	})
	if err != nil {
		hclog.Default().Info("Delete", err)
		return nil, err
	}
	return resp.R, err
}
