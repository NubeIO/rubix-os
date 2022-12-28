package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) SyncProducer(body *interfaces.SyncProducer) ([]*model.Consumer, error) {
	consumers, err := d.GetConsumers(api.Args{
		ProducerUUID:      nils.NewString(body.ProducerUUID),
		ProducerThingUUID: nils.NewString(body.ProducerThingUUID)})
	if err != nil {
		return nil, err
	}
	for _, consumer := range consumers {
		consumer.ProducerThingName = body.ProducerThingName
		_, _ = d.UpdateConsumer(consumer.UUID, consumer)
	}
	return consumers, nil
}
