package database

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/interfaces/connection"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/urls"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"gorm.io/gorm"
)

func (d *GormDatabase) GetStreamClones(args api.Args) ([]*model.StreamClone, error) {
	var streamClonesModel []*model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query.Find(&streamClonesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamClonesModel, nil
}

func (d *GormDatabase) GetStreamCloneByArg(args api.Args) (*model.StreamClone, error) {
	var streamClonesModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query.Find(&streamClonesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamClonesModel, nil
}

func GetOneStreamCloneByArgTransaction(db *gorm.DB, args api.Args) (*model.StreamClone, error) {
	var streamClonesModel *model.StreamClone
	query := buildStreamCloneQueryTransaction(db, args)
	query.First(&streamClonesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamClonesModel, nil
}

func (d *GormDatabase) GetOneStreamCloneByArg(args api.Args) (*model.StreamClone, error) {
	return GetOneStreamCloneByArgTransaction(d.DB, args)
}

func (d *GormDatabase) GetStreamClone(uuid string, args api.Args) (*model.StreamClone, error) {
	var streamCloneModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query = query.Where("uuid = ? ", uuid).First(&streamCloneModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamCloneModel, nil
}

func (d *GormDatabase) DeleteStreamClone(uuid string) (bool, error) {
	streamCloneModel, err := d.GetStreamClone(uuid, api.Args{})
	if err != nil {
		return false, err
	}
	if boolean.IsTrue(streamCloneModel.CreatedFromAutoMapping) {
		return false, errors.New("can't delete auto-mapped stream clone")
	}
	query := d.DB.Delete(&streamCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteOneStreamCloneByArgs(args api.Args) (bool, error) {
	var streamCloneModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	if err := query.First(&streamCloneModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&streamCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) updateStreamClone(uuid string, body *model.StreamClone) error {
	query := d.DB.Where("uuid = ?", uuid).Updates(body)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (d *GormDatabase) SyncStreamCloneConsumers(uuid string, args api.Args) ([]*interfaces.SyncModel, error) {
	streamClone, _ := d.GetStreamClone(uuid, api.Args{WithConsumers: true})
	if streamClone == nil {
		return nil, errors.New("no stream_clone")
	}
	flowNetworkClone, _ := d.GetFlowNetworkClone(streamClone.FlowNetworkCloneUUID, api.Args{})
	cli := client.NewFlowClientCliFromFNC(flowNetworkClone)
	var outputs []*interfaces.SyncModel
	localCli := client.NewLocalClient()
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	for _, consumer := range streamClone.Consumers {
		go d.syncConsumer(cli, localCli, consumer, args, channel)
	}
	for range streamClone.Consumers {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncConsumer(cli *client.FlowClient, localCli *client.FlowClient, consumer *model.Consumer,
	args api.Args, channel chan *interfaces.SyncModel) {
	rawProducer, err := cli.GetQueryMarshal(urls.SingularUrl(urls.ProducerUrl, consumer.ProducerUUID), model.Producer{})
	var output interfaces.SyncModel
	if err != nil {
		output = interfaces.SyncModel{UUID: consumer.UUID, IsError: true, Message: nstring.New(err.Error())}
		consumer.Connection = connection.Broken.String()
		consumer.Message = err.Error()
	} else {
		output = interfaces.SyncModel{UUID: consumer.UUID, IsError: false}
		producer := rawProducer.(*model.Producer)
		consumer.Connection = connection.Connected.String()
		consumer.Message = nstring.NotAvailable
		consumer.ProducerThingName = producer.ProducerThingName
		consumer.ProducerThingUUID = producer.ProducerThingUUID
		consumer.ProducerThingClass = producer.ProducerThingClass
		consumer.ProducerThingType = producer.ProducerThingType
	}
	// This is for syncing child descendants
	if args.WithWriters {
		url := urls.GetUrl(urls.ConsumerWritersSyncUrl, consumer.UUID)
		_, _ = localCli.GetQuery(url)
	}
	d.DB.Where("uuid = ?", consumer.UUID).Updates(consumer)
	channel <- &output
}
