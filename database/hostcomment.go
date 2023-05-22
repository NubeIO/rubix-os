package database

import (
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) CreateHostComment(body *model.HostComment) (*model.HostComment, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.HostComment)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateHostComment(uuid string, body *model.HostComment) (*model.HostComment, error) {
	hostModel := new(model.HostComment)
	err := d.DB.Where("uuid = ?", uuid).Find(&hostModel).Updates(body).Error
	if err != nil {
		return nil, err
	}
	return hostModel, nil
}

func (d *GormDatabase) DeleteHostComment(uuid string) (bool, error) {
	var hostModel *model.HostComment
	query := d.DB.Where("uuid = ? ", uuid).Delete(&hostModel)
	return d.deleteResponseBuilder(query)
}
