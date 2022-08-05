package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/config"
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
	if body.CurrentWriterUUID == nil {
		if err := d.DB.Model(&producerModel).Update("current_writer_uuid", nil).Error; err != nil {
			return nil, err
		}
	}
	return producerModel, nil
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
	createCOVHistory bool, currentWriterUUID *string) error {
	producerModelBody := new(model.Producer)
	producerModelBody.CurrentWriterUUID = currentWriterUUID
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &uuid})
	for _, producer := range producers {
		err := d.producerPointWrite(producer.UUID, priority, presentValue, producerModelBody, createCOVHistory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) producerPointWrite(uuid string, priority *map[string]*float64, presentValue *float64,
	producerModelBody *model.Producer, createCOVHistory bool) error {
	producerModel, err := d.UpdateProducer(uuid, producerModelBody)
	if err != nil {
		log.Errorf("producer: issue on update producer err: %v\n", err)
		return errors.New("issue on update producer")
	}
	syncCOV := model.SyncCOV{Priority: priority}
	err = d.TriggerCOVToWriterClone(producerModel, &syncCOV)
	if err != nil {
		return err
	}
	if boolean.IsTrue(config.Get().ProducerHistory.Enable) && createCOVHistory &&
		boolean.IsTrue(producerModel.EnableHistory) && checkHistoryCovType(string(producerModel.HistoryType)) {
		ph := new(model.ProducerHistory)
		ph.ProducerUUID = uuid
		ph.CurrentWriterUUID = producerModel.CurrentWriterUUID
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
	producerModel, err := d.UpdateProducer(uuid, producerModelBody)
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

func (d *GormDatabase) GetProducersForCreateInterval() ([]*interfaces.ProducerIntervalHistory, error) {
	var producerIntervalHistory []*interfaces.ProducerIntervalHistory
	query := fmt.Sprintf("SELECT p.uuid,p.producer_thing_class,p.history_interval,ph.timestamp AS timestamp,pt.present_value "+
		"FROM producers p "+
		"LEFT JOIN (SELECT producer_uuid, MAX(timestamp) AS timestamp FROM producer_histories GROUP BY producer_uuid) ph "+
		"ON p.uuid = ph.producer_uuid "+
		"INNER JOIN points pt "+
		"ON p.producer_thing_uuid = pt.uuid "+
		"WHERE p.enable_history = %v AND p.history_type != '%s' AND p.history_interval > %d", true, model.HistoryTypeCov, 0)
	if err := d.DB.Raw(query).Scan(&producerIntervalHistory).Error; err != nil {
		return nil, err
	}
	return producerIntervalHistory, nil
}
