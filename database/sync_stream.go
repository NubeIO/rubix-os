package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) SyncStream(body *model.SyncStream) (*model.StreamClone, error) {
	var flowNetworkClone *model.FlowNetworkClone
	if err := d.DB.Where("global_uuid = ?", body.GlobalUUID).First(&flowNetworkClone).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("we don't have flow_network_clone with global_uuid=%s", body.GlobalUUID))
	}
	mStream, err := json.Marshal(body.Stream)
	if err != nil {
		return nil, err
	}
	streamClone := model.StreamClone{}
	if err = json.Unmarshal(mStream, &streamClone); err != nil {
		return nil, err
	}
	streamClone.FlowNetworkCloneUUID = flowNetworkClone.UUID
	streamClone.SourceUUID = body.Stream.UUID
	var streamClonesModel []*model.StreamClone
	if err = d.DB.Where("source_uuid = ? ", streamClone.SourceUUID).Find(&streamClonesModel).Error; err != nil {
		return nil, err
	}
	if len(streamClonesModel) == 0 {
		// Stream clone suppose to delete with stream deletion. But when we do mapping changes and sync process
		// sometimes stream gets deleted whereas stream_clone doesn't get deleted, so when we sync again it shows the
		// conflict. And this block of code avoids such cases.
		if err = d.DB.
			Where("name = ? ", streamClone.Name).
			Where("created_from_auto_mapping IS TRUE").
			Find(&streamClonesModel).Error; err != nil {
			return nil, err
		}
		for _, scm := range streamClonesModel {
			d.DB.Delete(&scm)
		}
		streamClone.UUID = nuuid.MakeTopicUUID(model.CommonNaming.StreamClone)
		if err = d.DB.Create(&streamClone).Error; err != nil {
			return nil, err
		}
	} else {
		streamClone.UUID = streamClonesModel[0].UUID
		if err = d.DB.Model(&streamClonesModel[0]).Updates(streamClone).Error; err != nil {
			return nil, err
		}
	}
	return &streamClone, nil
}
