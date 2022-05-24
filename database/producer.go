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
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
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

func (d *GormDatabase) GetOneProducerByArgs(args api.Args) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.buildProducerQuery(args)
	if err := query.First(&producerModel).Error; err != nil {
		return nil, err
	}
	return producerModel, nil
}

func (d *GormDatabase) CreateProducer(body *model.Producer) (*model.Producer, error) {
	if body.ProducerThingUUID == "" {
		return nil, errors.New("please pass in a producer_thing_uuid i.e. uuid of that class")
	}
	if body.ProducerThingClass == "" {
		return nil, errors.New("please pass in a producer_thing_class i.e. point, job etc")
	}
	producerThingName := ""
	historyType, err := checkHistoryType(string(body.HistoryType))
	if err != nil {
		return nil, err
	}
	body.HistoryType = historyType

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
	_, err = d.GetStream(body.StreamUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Producer)
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = nuuid.MakeUUID()
	body.ProducerThingName = nameIsNil(producerThingName)
	if err = d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateProducer(uuid string, body *model.Producer) (*model.Producer, error) {
	var producerModel *model.Producer
	if err := d.DB.Where("uuid = ?", uuid).First(&producerModel).Error; err != nil {
		return nil, err
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&producerModel, body.Tags); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&producerModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return producerModel, nil
}

func (d *GormDatabase) DeleteProducer(uuid string) (bool, error) {
	producer, _ := d.GetProducer(uuid, api.Args{})
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
	for _, wc := range producer.WriterClones {
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
			_, err := cli.GetQueryMarshal(urls.SingularUrl(urls.WriterUrl, wc.SourceUUID), model.Writer{})
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
		outputs = append(outputs, &output)
	}
	return outputs, nil
}

type Point struct {
	Priority model.Priority `json:"priority"`
}

func (d *GormDatabase) ProducersPointWrite(uuid string, priority *map[string]*float64, presentValue *float64) error {
	producerModelBody := new(model.Producer)
	producerModelBody.CurrentWriterUUID = uuid // TODO: check current_writer_uuid is needed or not
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &uuid})
	for _, producer := range producers {
		err := d.producerPointWrite(producer.UUID, priority, presentValue, producerModelBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) producerPointWrite(uuid string, priority *map[string]*float64, presentValue *float64, producerModelBody *model.Producer) error {
	producerModel, err := d.UpdateProducer(uuid, producerModelBody) // TODO: check current_writer_uuid is needed or not
	if err != nil {
		log.Errorf("producer: issue on update producer err: %v\n", err)
		return errors.New("issue on update producer")
	}
	syncCOV := model.SyncCOV{Priority: priority}
	err = d.TriggerCOVToWriterClone(producerModel, &syncCOV)
	if err != nil {
		return err
	}
	if boolean.IsTrue(producerModel.EnableHistory) && checkHistoryCovType(string(producerModel.HistoryType)) {
		ph := new(model.ProducerHistory)
		ph.ProducerUUID = uuid
		ph.PresentValue = presentValue
		ph.Timestamp = time.Now().UTC()
		_, err = d.CreateProducerHistory(ph)
		if err != nil {
			log.Errorf("producer: issue on write history ProducerWriteHist: %v\n", err)
			return errors.New("issue on write history for point")
		}
	}
	return nil
}

func (d *GormDatabase) ProducersScheduleWrite(uuid string, body *model.ScheduleData) error {
	producerModelBody := new(model.Producer)
	producerModelBody.CurrentWriterUUID = uuid // TODO: check current_writer_uuid is needed or not
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
	producerModel, err := d.UpdateProducer(uuid, producerModelBody) // TODO: check current_writer_uuid is needed or not
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
		if wc.FlowFrameworkUUID == "" || fn.UUID == wc.FlowFrameworkUUID {
			cli := client.NewFlowClientCliFromFN(fn)
			_ = cli.SyncCOV(wc.SourceUUID, body)
		}
	}
	return nil
}

func (d *GormDatabase) GetProducersForCreateInterval() ([]*model.Producer, error) {
	var historyModel []*model.Producer
	query := d.DB.Where("enable_history = ? AND history_type != ?", true, model.HistoryTypeCov).
		Find(&historyModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyModel, nil
}
