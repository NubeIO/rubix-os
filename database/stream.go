package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) GetStreams(args api.Args) ([]*model.Stream, error) {
	var streamsModel []*model.Stream
	query := d.buildStreamQuery(args)
	query.Find(&streamsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamsModel, nil
}

func (d *GormDatabase) GetStream(uuid string, args api.Args) (*model.Stream, error) {
	var streamModel *model.Stream
	query := d.buildStreamQuery(args)
	query = query.Where("uuid = ? ", uuid).First(&streamModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamModel, nil
}

func (d *GormDatabase) GetStreamByArgs(args api.Args) ([]*model.Stream, error) {
	var streamsModel []*model.Stream
	query := d.buildStreamQuery(args)
	if err := query.Find(&streamsModel).Error; err != nil {
		return nil, err
	}
	return streamsModel, nil
}

func (d *GormDatabase) GetOneStreamByArgs(args api.Args) (*model.Stream, error) {
	var streamModel *model.Stream
	query := d.buildStreamQuery(args)
	if err := query.First(&streamModel).Error; err != nil {
		return nil, err
	}
	return streamModel, nil
}

func (d *GormDatabase) CreateStream(body *model.Stream) (*model.Stream, error) {
	stream, _ := d.GetOneStreamByArgs(api.Args{Name: nils.NewString(body.Name)})
	if stream != nil {
		return stream, errors.New("an existing stream with this name exists")
	}
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Stream)
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = nuuid.MakeUUID()
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	_ = d.syncAfterCreateUpdateStream(body)
	return body, nil
}

func (d *GormDatabase) GetFlowNetworksFromStreamUUID(streamUUID string) (*[]model.FlowNetwork, error) {
	var flowNetworks *[]model.FlowNetwork
	err := d.DB.
		Joins("JOIN flow_networks_streams ON flow_networks_streams.flow_network_uuid = flow_networks.uuid").
		Where("flow_networks_streams.stream_uuid IN (?)", streamUUID).
		Find(&flowNetworks).
		Error
	if err != nil {
		return nil, err
	}
	return flowNetworks, nil
}

func (d *GormDatabase) UpdateStream(uuid string, body *model.Stream, checkAm bool) (*model.Stream, error) {
	var streamModel *model.Stream
	if err := d.DB.Preload("FlowNetworks").Where("uuid = ?", uuid).First(&streamModel).Error; err != nil {
		return nil, err
	}
	if boolean.IsTrue(streamModel.CreatedFromAutoMapping) && checkAm {
		return nil, errors.New("can't update auto-mapped stream")
	}
	if len(body.FlowNetworks) > 0 {
		if err := d.DB.Model(&streamModel).Association("FlowNetworks").Replace(body.FlowNetworks); err != nil {
			return nil, err
		}
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&streamModel, body.Tags); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&streamModel).Updates(body).Error; err != nil {
		return nil, err
	}
	_ = d.syncAfterCreateUpdateStream(streamModel)
	return streamModel, nil
}

func (d *GormDatabase) DeleteStream(uuid string) (bool, error) {
	streamModel, err := d.GetStream(uuid, api.Args{WithFlowNetworks: true})
	if err != nil {
		return false, err
	}
	if boolean.IsTrue(streamModel.CreatedFromAutoMapping) {
		return false, errors.New("can't delete auto-mapped stream")
	}
	var aType = api.ArgsType
	for _, fn := range streamModel.FlowNetworks {
		cli := client.NewFlowClientCliFromFN(fn)
		url := urls.SingularUrlByArg(urls.StreamCloneUrl, aType.SourceUUID, streamModel.UUID)
		_ = cli.DeleteQuery(url)
	}
	query := d.DB.Delete(&streamModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) SyncStreamProducers(uuid string, args api.Args) ([]*interfaces.SyncModel, error) {
	stream, _ := d.GetStream(uuid, api.Args{WithProducers: true})
	localCli := client.NewLocalClient()
	channel := make(chan bool)
	defer close(channel)
	// This is for syncing child descendants
	if args.WithWriterClones {
		for _, producer := range stream.Producers {
			go d.syncProducer(localCli, producer, channel)
		}
		for range stream.Producers {
			<-channel
		}
	}
	return nil, nil
}

func (d *GormDatabase) syncProducer(localCli *client.FlowClient, producer *model.Producer, channel chan bool) {
	url := urls.GetUrl(urls.ProducerWriterClonesSyncUrl, producer.UUID)
	_, _ = localCli.GetQuery(url)
	channel <- true
}

func (d *GormDatabase) syncAfterCreateUpdateStream(body *model.Stream) error {
	flowNetworks, err := d.GetFlowNetworksFromStreamUUID(body.UUID)
	if err != nil {
		return err
	} else if len(*flowNetworks) == 0 {
		return nil
	}
	deviceInfo, err := deviceinfo.GetDeviceInfo()
	if err != nil {
		return err
	}
	for _, fn := range *flowNetworks {
		_ = d.SyncStreamFunction(&fn, body, deviceInfo)
	}
	return nil
}

func (d *GormDatabase) SyncStreamFunction(fn *model.FlowNetwork, stream *model.Stream, deviceInfo *model.DeviceInfo) error {
	cli := client.NewFlowClientCliFromFN(fn)
	syncStreamBody := model.SyncStream{
		GlobalUUID: deviceInfo.GlobalUUID,
		Stream:     stream,
	}
	_, err := cli.SyncStream(&syncStreamBody)
	if err != nil {
		return err
	}
	return nil
}
