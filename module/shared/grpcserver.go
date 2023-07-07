package shared

import (
	"context"
	"errors"
	"github.com/NubeIO/rubix-os/module/proto"
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
	err = m.Impl.Init(dbHelper, req.ModuleName)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (m *GRPCServer) Enable(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	log.Debug("gRPC Enable server has been called...")
	err := m.Impl.Enable()
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (m *GRPCServer) Disable(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	log.Debug("gRPC Disable server has been called...")
	err := m.Impl.Disable()
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (m *GRPCServer) ValidateAndSetConfig(ctx context.Context, req *proto.ConfigBody) (*proto.Response, error) {
	log.Debug("gRPC Disable server has been called...")
	bytes, err := m.Impl.ValidateAndSetConfig(req.Config)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: bytes}, nil
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
	}, nil
}

func (m *GRPCServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.Response, error) {
	log.Debug("gRPC Get server has been called...")
	r, err := m.Impl.Get(req.Path)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, nil
}

func (m *GRPCServer) Post(ctx context.Context, req *proto.PostRequest) (*proto.Response, error) {
	log.Debug("gRPC Post server has been called...")
	r, err := m.Impl.Post(req.Path, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, nil
}

func (m *GRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Response, error) {
	log.Debug("gRPC Put server has been called...")
	r, err := m.Impl.Put(req.Path, req.Uuid, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, nil
}

func (m *GRPCServer) Patch(ctx context.Context, req *proto.PatchRequest) (*proto.Response, error) {
	log.Debug("gRPC Patch server has been called...")
	r, err := m.Impl.Patch(req.Path, req.Uuid, req.Body)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, nil
}

func (m *GRPCServer) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.Response, error) {
	log.Debug("gRPC Delete server has been called...")
	r, err := m.Impl.Delete(req.Path, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &proto.Response{R: r}, nil
}

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCDBHelperClient struct{ client proto.DBHelperClient }

func (m *GRPCDBHelperClient) GetWithoutParam(path, args string) ([]byte, error) {
	resp, err := m.client.GetWithoutParam(context.Background(), &proto.GetWithoutParamRequest{
		Path: path,
		Args: args,
	})
	if err != nil {
		log.Error("GetList: ", err)
		return nil, err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("GetList: ", errStr)
		return nil, errors.New(errStr)
	}
	return resp.R, nil
}

func (m *GRPCDBHelperClient) Get(path, uuid, args string) ([]byte, error) {
	resp, err := m.client.Get(context.Background(), &proto.GetRequest{
		Path: path,
		Uuid: uuid,
		Args: args,
	})
	if err != nil {
		log.Error("Get: ", err)
		return nil, err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("Get: ", errStr)
		return nil, errors.New(errStr)
	}
	return resp.R, nil
}

func (m *GRPCDBHelperClient) Post(path string, body []byte) ([]byte, error) {
	resp, err := m.client.Post(context.Background(), &proto.PostRequest{
		Path: path,
		Body: body,
	})
	if err != nil {
		log.Error("Post: ", err)
		return nil, err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("Post: ", errStr)
		return nil, errors.New(errStr)
	}
	return resp.R, nil
}

func (m *GRPCDBHelperClient) Put(path, uuid string, body []byte) ([]byte, error) {
	resp, err := m.client.Put(context.Background(), &proto.PutRequest{
		Path: path,
		Uuid: uuid,
		Body: body,
	})
	if err != nil {
		log.Error("Put: ", err)
		return nil, err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("Put: ", errStr)
		return nil, errors.New(errStr)
	}
	return resp.R, nil
}

func (m *GRPCDBHelperClient) Patch(path, uuid string, body []byte) ([]byte, error) {
	resp, err := m.client.Patch(context.Background(), &proto.PatchRequest{
		Path: path,
		Uuid: uuid,
		Body: body,
	})
	if err != nil {
		log.Error("Patch: ", err)
		return nil, err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("Patch: ", errStr)
		return nil, errors.New(errStr)
	}
	return resp.R, nil
}

func (m *GRPCDBHelperClient) Delete(path, uuid string) ([]byte, error) {
	resp, err := m.client.Delete(context.Background(), &proto.DeleteRequest{
		Path: path,
		Uuid: uuid,
	})
	if err != nil {
		log.Error("Delete: ", err)
		return nil, err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("Delete: ", errStr)
		return nil, errors.New(errStr)
	}
	return resp.R, nil
}

func (m *GRPCDBHelperClient) SetErrorsForAll(path, uuid, message, messageLevel, messageCode string, doPoints bool) error {
	resp, err := m.client.SetErrorsForAll(context.Background(), &proto.SetErrorsForAllRequest{
		Path:         path,
		Uuid:         uuid,
		Message:      message,
		MessageLevel: messageLevel,
		MessageCode:  messageCode,
		DoPoints:     doPoints,
	})
	if err != nil {
		log.Error("SetErrorsForAll: ", err)
		return err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("SetErrorsForAll: ", errStr)
		return errors.New(errStr)
	}
	return nil
}

func (m *GRPCDBHelperClient) ClearErrorsForAll(path, uuid string, doPoints bool) error {
	resp, err := m.client.ClearErrorsForAll(context.Background(), &proto.ClearErrorsForAllRequest{
		Path:     path,
		Uuid:     uuid,
		DoPoints: doPoints,
	})
	if err != nil {
		log.Error("ClearErrorsForAll: ", err)
		return err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("ClearErrorsForAll: ", errStr)
		return errors.New(errStr)
	}
	return nil
}

func (m *GRPCDBHelperClient) WizardNewNetworkDevicePoint(plugin string, network, device, point []byte) (bool, error) {
	resp, err := m.client.WizardNewNetworkDevicePoint(context.Background(), &proto.WizardNewNetworkDevicePointRequest{
		Plugin:  plugin,
		Network: network,
		Device:  device,
		Point:   point,
	})
	if err != nil {
		log.Error("WizardNewNetworkDevicePoint: ", err)
		return false, err
	}
	if resp.E != nil {
		errStr := string(resp.E)
		log.Error("WizardNewNetworkDevicePoint: ", errStr)
		return false, errors.New(errStr)
	}
	return true, nil
}
