package database

import (
	"encoding/json"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

func (d *GormDatabase) SyncStream(body *model.StreamSync) (*model.StreamClone, error) {
	var flowNetworkClone *model.FlowNetworkClone
	if err := d.DB.Find(&flowNetworkClone).Error; err != nil {
		return nil, err
	}
	mStream, err := json.Marshal(body.Stream)
	if err != nil {
		return nil, err
	}
	streamClone := model.StreamClone{}
	if err = json.Unmarshal(mStream, &streamClone); err != nil {
		return nil, err
	}
	streamClone.UUID = utils.MakeTopicUUID(model.CommonNaming.WriterClone)
	streamClone.FlowNetworkCloneUUID = flowNetworkClone.UUID
	streamClone.SourceUUID = body.Stream.UUID
	var streamClonesModel []*model.StreamClone
	if err = d.DB.Where("source_uuid = ? ", streamClone.SourceUUID).Find(&streamClonesModel).Error; err != nil {
		return nil, err
	}
	if len(streamClonesModel) == 0 {
		if err = d.DB.Create(streamClone).Error; err != nil {
			return nil, err
		}
	} else {
		if err = d.DB.Model(&streamClonesModel[0]).Updates(streamClone).Error; err != nil {
			return nil, err
		}
	}
	return &streamClone, nil
}
