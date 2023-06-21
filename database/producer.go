package database

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/interfaces/connection"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/urls"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Producer struct {
	*model.Producer
}

func (d *GormDatabase) GetProducers(args api.Args) ([]*model.Producer, error) {
	var producersModel []*model.Producer
	query := d.buildProducerQuery(args)
	if err := query.Find(&producersModel).Error; err != nil {
		return nil, err
	}
	return producersModel, nil
}

func (d *GormDatabase) GetProducer(uuid string, args api.Args) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.buildProducerQuery(args)
	if err := query.Where("uuid = ?", uuid).First(&producerModel).Error; err != nil {
		return nil, err
	}
	return producerModel, nil
}

func (d *GormDatabase) GetOneProducerByArgsTransaction(db *gorm.DB, args api.Args) (*model.Producer, error) {
	var producerModel *model.Producer
	query := buildProducerQueryTransaction(db, args)
	if err := query.First(&producerModel).Error; err != nil {
		return nil, err
	}
	return producerModel, nil
}

func (d *GormDatabase) GetOneProducerByArgs(args api.Args) (*model.Producer, error) {
	return d.GetOneProducerByArgsTransaction(d.DB, args)
}

func (d *GormDatabase) CreateProducer(body *model.Producer) (*model.Producer, error) {
	if body.ProducerThingUUID == "" {
		return nil, errors.New("please pass in a producer_thing_uuid i.e. uuid of that class")
	}
	if body.ProducerThingClass == "" {
		return nil, errors.New("please pass in a producer_thing_class i.e. point, job etc")
	}
	producerThingName := ""

	switch body.ProducerThingClass {
	case model.ThingClass.Point:
		pnt, err := d.GetPoint(body.ProducerThingUUID, api.Args{})
		if err != nil {
			return nil, errors.New("point not found, please supply a valid point producer_thing_uuid")
		}
		producerThingName = pnt.Name
	case model.ThingClass.Schedule:
		sch, err := d.GetSchedule(body.ProducerThingUUID)
		if err != nil {
			return nil, errors.New("schedule not found, please supply a valid point producer_thing_uuid")
		}
		producerThingName = sch.Name
	default:
		return nil, errors.New("we are not supporting producer_thing_class other than point & schedule")
	}
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Producer)
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = nuuid.MakeUUID()
	body.ProducerThingName = nameIsNil(producerThingName)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateProducer(uuid string, body *model.Producer, checkAutoMap bool) (*model.Producer, error) {
	var producerModel *model.Producer
	if err := d.DB.Where("uuid = ?", uuid).First(&producerModel).Error; err != nil {
		return nil, err
	}
	if err := d.updateTags(&producerModel, body.Tags); err != nil {
		return nil, err
	}
	syncConsumer := body.ProducerThingName != "" && body.ProducerThingName != producerModel.ProducerThingName
	if err := d.DB.Model(&producerModel).Updates(body).Error; err != nil {
		return nil, err
	}
	if body.CurrentWriterUUID == nil {
		if err := d.DB.Model(&producerModel).Update("current_writer_uuid", nil).Error; err != nil {
			return nil, err
		}
	}
	if syncConsumer {
		stream, _ := d.GetStream(producerModel.StreamUUID, api.Args{WithFlowNetworks: true})
		syncBody := interfaces.SyncProducer{
			ProducerUUID:      producerModel.UUID,
			ProducerThingUUID: producerModel.ProducerThingUUID,
			ProducerThingName: producerModel.ProducerThingName,
		}
		for _, fn := range stream.FlowNetworks {
			cli := client.NewFlowClientCliFromFN(fn)
			_, _ = cli.SyncProducer(&syncBody)
		}
	}
	return producerModel, nil
}

func (d *GormDatabase) UpdateProducerByProducerThingUUID(producerThingUUID string, producerThingName string) {
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: nils.NewString(producerThingUUID)})
	for _, producer := range producers {
		if producer.ProducerThingName != producerThingName {
			producer.ProducerThingName = producerThingName
			go d.UpdateProducer(producer.UUID, producer, false)
		}
	}
}

func (d *GormDatabase) DeleteProducer(uuid string) (bool, error) {
	producer, err := d.GetProducer(uuid, api.Args{})
	if err != nil {
		return false, err
	}
	stream, _ := d.GetStream(producer.StreamUUID, api.Args{WithFlowNetworks: true})
	aType := api.ArgsType
	url := urls.PluralUrlByArg(urls.ConsumerUrl, aType.ProducerUUID, producer.UUID)
	for _, fn := range stream.FlowNetworks {
		cli := client.NewFlowClientCliFromFN(fn)
		_ = cli.DeleteQuery(url)
	}
	query := d.DB.Delete(&producer)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) SyncProducerWriterClones(uuid string) ([]*interfaces.SyncModel, error) {
	producer, _ := d.GetProducer(uuid, api.Args{WithWriterClones: true})
	if producer == nil {
		return nil, errors.New("producer not found")
	}
	stream, _ := d.GetStream(producer.StreamUUID, api.Args{WithFlowNetworks: true})
	var outputs []*interfaces.SyncModel
	flowNetworkMap := map[string]*model.FlowNetwork{}
	for _, fn := range stream.FlowNetworks {
		flowNetworkMap[fn.UUID] = fn
	}
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	for _, wc := range producer.WriterClones {
		go d.syncWriterClone(flowNetworkMap, wc, channel)
	}
	for range producer.WriterClones {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncWriterClone(flowNetworkMap map[string]*model.FlowNetwork, wc *model.WriterClone,
	channel chan *interfaces.SyncModel) {
	var output interfaces.SyncModel
	fn, exists := flowNetworkMap[wc.FlowFrameworkUUID]
	wc.Connection = connection.Connected.String()
	wc.Message = nstring.NotAvailable
	if !exists {
		msg := "FlowNetwork is broken!"
		output = interfaces.SyncModel{
			UUID:    wc.UUID,
			IsError: true,
			Message: nstring.New(msg),
		}
		wc.Connection = connection.Broken.String()
		wc.Message = msg
	} else {
		cli := client.NewFlowClientCliFromFN(fn)
		_, err := cli.GetQuery(urls.SingularUrl(urls.WriterUrl, wc.SourceUUID))
		if err != nil {
			output = interfaces.SyncModel{
				UUID:    wc.UUID,
				IsError: true,
				Message: nstring.New(err.Error()),
			}
			wc.Connection = connection.Broken.String()
			wc.Message = err.Error()
		}
	}
	_ = d.updateWriterClone(wc.UUID, wc)
	channel <- &output
}

type Point struct {
	Priority model.Priority `json:"priority"`
}

func (d *GormDatabase) ProducersPointWrite(uuid string, priority *map[string]*float64, presentValue *float64,
	currentWriterUUID *string) error {
	producerModelBody := new(model.Producer)
	producerModelBody.CurrentWriterUUID = currentWriterUUID
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &uuid, Enable: boolean.NewTrue()})
	for _, producer := range producers {
		err := d.producerPointWrite(producer.UUID, priority, presentValue, producerModelBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) producerPointWrite(uuid string, priority *map[string]*float64, presentValue *float64,
	producerModelBody *model.Producer) error {
	producerModel, err := d.UpdateProducer(uuid, producerModelBody, false)
	if err != nil {
		log.Errorf("producer: issue on update producer err: %v\n", err)
		return errors.New("issue on update producer")
	}
	syncCOV := model.SyncCOV{Priority: priority, PresentValue: presentValue}
	err = d.TriggerCOVToWriterClone(producerModel, &syncCOV)
	if err != nil {
		return err
	}
	return nil
}

func (d *GormDatabase) ProducersScheduleWrite(uuid string, body *model.ScheduleData) error {
	producerModelBody := new(model.Producer)
	producerModelBody.CurrentWriterUUID = nil
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &uuid})
	for _, producer := range producers {
		err := d.producerScheduleWrite(producer.UUID, body, producerModelBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) producerScheduleWrite(uuid string, scheduleData *model.ScheduleData, producerModelBody *model.Producer) error {
	producerModel, err := d.UpdateProducer(uuid, producerModelBody, false)
	if err != nil {
		log.Errorf("producer: issue on update producer err: %v\n", err)
		return errors.New("issue on update producer")
	}

	syncCOV := model.SyncCOV{Schedule: scheduleData}
	err = d.TriggerCOVToWriterClone(producerModel, &syncCOV)
	if err != nil {
		return err
	}
	return nil
}

func (d *GormDatabase) TriggerCOVToWriterClone(producer *model.Producer, body *model.SyncCOV) error {
	wcs, err := d.GetWriterClones(api.Args{ProducerUUID: nstring.NewStringAddress(producer.UUID)})
	if err != nil {
		return errors.New("error on getting writer clones from producer_uuid")
	}
	for _, wc := range wcs {
		_ = d.TriggerCOVFromWriterCloneToWriter(producer, wc, body)
	}
	return nil
}

func (d *GormDatabase) TriggerCOVFromWriterCloneToWriter(producer *model.Producer, wc *model.WriterClone, body *model.SyncCOV) error {
	stream, _ := d.GetStream(producer.StreamUUID, api.Args{WithFlowNetworks: true})
	for _, fn := range stream.FlowNetworks {
		// TODO: wc.FlowFrameworkUUID == "" remove from condition; it's here coz old deployment doesn't used to have that value
		if (wc.FlowFrameworkUUID == "" || fn.UUID == wc.FlowFrameworkUUID) && boolean.IsTrue(stream.Enable) {
			cli := client.NewFlowClientCliFromFN(fn)
			_ = cli.SyncCOV(wc.SourceUUID, body)
		}
	}
	return nil
}
