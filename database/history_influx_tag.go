package database

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

func (d *GormDatabase) GetHistoryInfluxTags(producerUuid string) ([]*model.HistoryInfluxTag, error) {
	var influxHistoryTags []*model.HistoryInfluxTag
	return influxHistoryTags, nil
}
