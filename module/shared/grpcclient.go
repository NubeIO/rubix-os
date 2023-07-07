package shared

import (
	"context"
	"github.com/NubeIO/rubix-os/module/proto"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCClient is an implementation of Module that talks over RPC.
type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.ModuleClient
}

func (m *GRPCClient) Init(dbHelper DBHelper, moduleName string) error {
	log.Debug("gRPC Init client has been called...")
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
		AddServer:  brokerID,
		ModuleName: moduleName,
	})

	// s.Stop() // TODO: we haven't closed this
	return err
}

func (m *GRPCClient) Enable() error {
	log.Debug("gRPC Enable client has been called...")
	_, err := m.client.Enable(context.Background(), &proto.Empty{})
	return err
}

func (m *GRPCClient) Disable() error {
	log.Debug("gRPC Disable client has been called...")
	_, err := m.client.Disable(context.Background(), &proto.Empty{})
	return err
}

func (m *GRPCClient) ValidateAndSetConfig(config []byte) ([]byte, error) {
	log.Debug("gRPC ValidateAndSetConfig client has been called...")
	resp, err := m.client.ValidateAndSetConfig(context.Background(), &proto.ConfigBody{Config: config})
	if err != nil {
		return nil, err
	}
	return resp.R, nil
}

func (m *GRPCClient) GetInfo() (*Info, error) {
	log.Debug("gRPC GetInfo client has been called...")
	resp, err := m.client.GetInfo(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}
	return &Info{
		Name:       resp.Name,
		Author:     resp.Author,
		Website:    resp.Website,
		License:    resp.License,
		HasNetwork: resp.HasNetwork,
	}, nil
}

func (m *GRPCClient) Get(path string) ([]byte, error) {
	log.Debug("gRPC Get client has been called...")
	resp, err := m.client.Get(context.Background(), &proto.GetRequest{
		Path: path,
	})
	if err != nil {
		return nil, err
	}
	return resp.R, nil
}

func (m *GRPCClient) Post(path string, body []byte) ([]byte, error) {
	log.Debug("gRPC Post client has been called...")
	resp, err := m.client.Post(context.Background(), &proto.PostRequest{
		Path: path,
		Body: body,
	})
	if err != nil {
		return nil, err
	}
	return resp.R, nil
}

func (m *GRPCClient) Put(path, uuid string, body []byte) ([]byte, error) {
	log.Debugf("gRPC Put client has been called...")
	resp, err := m.client.Put(context.Background(), &proto.PutRequest{
		Path: path,
		Uuid: uuid,
		Body: body,
	})
	if err != nil {
		return nil, err
	}
	return resp.R, nil
}

func (m *GRPCClient) Patch(path, uuid string, body []byte) ([]byte, error) {
	log.Debug("gRPC Patch client has been called...")
	resp, err := m.client.Patch(context.Background(), &proto.PatchRequest{
		Path: path,
		Uuid: uuid,
		Body: body,
	})
	if err != nil {
		return nil, err
	}
	return resp.R, nil
}

func (m *GRPCClient) Delete(path, uuid string) ([]byte, error) {
	log.Debug("gRPC Delete client has been called...")
	resp, err := m.client.Delete(context.Background(), &proto.DeleteRequest{
		Path: path,
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return resp.R, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCDBHelperServer struct {
	// This is the real implementation
	Impl DBHelper
}

func (m *GRPCDBHelperServer) GetWithoutParam(ctx context.Context, req *proto.GetWithoutParamRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.GetWithoutParam(req.Path, req.Args)
	if err != nil {
		return &proto.Response{R: nil, E: []byte(err.Error())}, nil
	}
	return &proto.Response{R: r, E: nil}, nil
}

func (m *GRPCDBHelperServer) Get(ctx context.Context, req *proto.GetRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Get(req.Path, req.Uuid, req.Args)
	if err != nil {
		return &proto.Response{R: nil, E: []byte(err.Error())}, nil
	}
	return &proto.Response{R: r, E: nil}, nil
}

func (m *GRPCDBHelperServer) Post(ctx context.Context, req *proto.PostRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Post(req.Path, req.Body)
	if err != nil {
		return &proto.Response{R: nil, E: []byte(err.Error())}, nil
	}
	return &proto.Response{R: r, E: nil}, nil
}

func (m *GRPCDBHelperServer) Put(ctx context.Context, req *proto.PutRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Put(req.Path, req.Uuid, req.Body)
	if err != nil {
		return &proto.Response{R: nil, E: []byte(err.Error())}, nil
	}
	return &proto.Response{R: r, E: nil}, nil
}

func (m *GRPCDBHelperServer) Patch(ctx context.Context, req *proto.PatchRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Patch(req.Path, req.Uuid, req.Body)
	if err != nil {
		return &proto.Response{R: nil, E: []byte(err.Error())}, nil
	}
	return &proto.Response{R: r, E: nil}, nil
}

func (m *GRPCDBHelperServer) Delete(ctx context.Context, req *proto.DeleteRequest) (resp *proto.Response, err error) {
	r, err := m.Impl.Delete(req.Path, req.Uuid)
	if err != nil {
		return &proto.Response{R: nil, E: []byte(err.Error())}, nil
	}
	return &proto.Response{R: r, E: nil}, nil
}

func (m *GRPCDBHelperServer) SetErrorsForAll(ctx context.Context, request *proto.SetErrorsForAllRequest) (*proto.ErrorResponse, error) {
	err := m.Impl.SetErrorsForAll(
		request.Path,
		request.Uuid,
		request.Message,
		request.MessageLevel,
		request.MessageCode,
		request.DoPoints,
	)
	if err != nil {
		return &proto.ErrorResponse{E: []byte(err.Error())}, nil
	}
	return &proto.ErrorResponse{E: nil}, nil
}

func (m *GRPCDBHelperServer) ClearErrorsForAll(ctx context.Context, request *proto.ClearErrorsForAllRequest) (*proto.ErrorResponse, error) {
	err := m.Impl.ClearErrorsForAll(request.Path, request.Uuid, request.DoPoints)
	if err != nil {
		return &proto.ErrorResponse{E: []byte(err.Error())}, nil
	}
	return &proto.ErrorResponse{E: nil}, nil
}

func (m *GRPCDBHelperServer) WizardNewNetworkDevicePoint(ctx context.Context, request *proto.WizardNewNetworkDevicePointRequest) (*proto.BoolResponse, error) {
	_, err := m.Impl.WizardNewNetworkDevicePoint(request.Plugin, request.Network, request.Device, request.Point)
	if err != nil {
		return &proto.BoolResponse{R: false, E: []byte(err.Error())}, nil
	}
	return &proto.BoolResponse{R: true, E: nil}, nil
}
