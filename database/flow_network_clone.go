package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) GetFlowNetworkClones(args api.Args) ([]*model.FlowNetworkClone, error) {
	var flowNetworkClonesModel []*model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Find(&flowNetworkClonesModel).Error; err != nil {
		return nil, err
	}
	return flowNetworkClonesModel, nil
}

func (d *GormDatabase) GetFlowNetworkClone(uuid string, args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&flowNetworkCloneModel).Error; err != nil {
		return nil, err

	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) DeleteFlowNetworkClone(uuid string) (bool, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.DB.Where("uuid = ? ", uuid).Delete(&flowNetworkCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) GetOneFlowNetworkCloneByArgs(args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.First(&flowNetworkCloneModel).Error; err != nil {
		return nil, err
	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) DeleteOneFlowNetworkCloneByArgs(args api.Args) (bool, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.First(&flowNetworkCloneModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&flowNetworkCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) RefreshFlowNetworkClonesConnections() (*bool, error) {
	var flowNetworkClones []*model.FlowNetworkClone
	d.DB.Where("is_master_slave IS NOT TRUE AND is_remote IS TRUE").Find(&flowNetworkClones)
	for _, fnc := range flowNetworkClones {
		accessToken, err := client.GetFlowToken(*fnc.FlowIP, *fnc.FlowPort, *fnc.FlowUsername, *fnc.FlowPassword)
		fncModel := model.FlowNetworkClone{}
		if err != nil {
			fncModel.IsError = boolean.NewTrue()
			fncModel.ErrorMsg = nstring.NewStringAddress(err.Error())
			fncModel.FlowToken = fnc.FlowToken
		} else {
			fncModel.IsError = boolean.NewFalse()
			fncModel.ErrorMsg = nil
			fncModel.FlowToken = accessToken
		}
		// here `.Select` is needed because NULL value needs to set on is_error=false
		if err := d.DB.Model(&fnc).Select("IsError", "ErrorMsg", "FlowToken").Updates(fncModel).Error; err != nil {
			log.Error(err)
		}
	}
	return boolean.NewTrue(), nil
}

func (d *GormDatabase) SyncFlowNetworkClones(args api.Args) ([]*interfaces.SyncModel, error) {
	fncs, _ := d.GetFlowNetworkClones(api.Args{})
	var outputs []*interfaces.SyncModel
	params := urls.GenerateFNCUrlParams(args)
	localCli := client.NewLocalClient()
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	for _, fnc := range fncs {
		go d.syncFlowNetworkClone(localCli, fnc, args, params, channel)
	}
	for range fncs {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncFlowNetworkClone(localCli *client.FlowClient, fnc *model.FlowNetworkClone, args api.Args,
	params string, channel chan *interfaces.SyncModel) {
	var output interfaces.SyncModel
	cli := client.NewFlowClientCliFromFNC(fnc)
	_, err := cli.GetQuery(urls.SingularUrl(urls.FlowNetworkUrl, fnc.SourceUUID))
	if err != nil {
		fnc.Message = err.Error()
		fnc.Connection = connection.Broken.String()
		output = interfaces.SyncModel{UUID: fnc.UUID, IsError: true, Message: nstring.New(err.Error())}
	} else {
		fnc.Message = nstring.NotAvailable
		fnc.Connection = connection.Connected.String()
		output = interfaces.SyncModel{UUID: fnc.UUID, IsError: false}
	}
	// This is for syncing child descendants
	if args.WithStreamClones == true {
		url := urls.GetUrl(urls.FlowNetworkCloneStreamClonesSyncUrl, fnc.UUID) + params
		_, _ = localCli.GetQuery(url)
	}
	d.DB.Where("uuid = ?", fnc.UUID).Updates(fnc)
	channel <- &output
}

func (d *GormDatabase) SyncFlowNetworkCloneStreamClones(uuid string, args api.Args) ([]*interfaces.SyncModel, error) {
	var outputs []*interfaces.SyncModel
	fnc, _ := d.GetFlowNetworkClone(uuid, api.Args{WithStreamClones: true})
	if fnc == nil {
		return nil, errors.New("no flow_network_clone")
	}
	params := urls.GenerateFNCUrlParams(args)
	localCli := client.NewLocalClient()
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	for _, sc := range fnc.StreamClones {
		go d.syncStreamClone(localCli, fnc, sc, args, params, channel)
	}
	for range fnc.StreamClones {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncStreamClone(localCli *client.FlowClient, fnc *model.FlowNetworkClone, sc *model.StreamClone,
	args api.Args, params string, channel chan *interfaces.SyncModel) {
	var output interfaces.SyncModel
	cli := client.NewFlowClientCliFromFNC(fnc)
	_, err := cli.GetQuery(urls.SingularUrl(urls.StreamUrl, sc.SourceUUID))
	if err != nil {
		output = interfaces.SyncModel{
			UUID:    sc.UUID,
			IsError: true,
			Message: nstring.New(err.Error()),
		}
		sc.Connection = connection.Broken.String()
		sc.Message = err.Error()
	} else {
		output = interfaces.SyncModel{
			UUID:    sc.UUID,
			IsError: false,
		}
		sc.Connection = connection.Connected.String()
		sc.Message = nstring.NotAvailable
	}
	_ = d.updateStreamClone(sc.UUID, sc)
	// This is for syncing child descendants
	if args.WithConsumers == true {
		url := urls.GetUrl(urls.StreamCloneConsumersSyncUrl, sc.UUID) + params
		_, _ = localCli.GetQuery(url)
	}
	channel <- &output
}
