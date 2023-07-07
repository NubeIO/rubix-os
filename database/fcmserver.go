package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"github.com/NubeIO/rubix-os/utils/security"
)

func (d *GormDatabase) GetFcmServer() (*model.FcmServer, error) {
	var fcmServerModel *model.FcmServer
	if err := d.DB.First(&fcmServerModel).Error; err != nil {
		return nil, err
	}
	fcmServerModel.Key = security.Decrypt(fcmServerModel.Key)
	return fcmServerModel, nil
}

func (d *GormDatabase) GetFcmServerKey() string {
	fcmServerModel, _ := d.GetFcmServer()
	if fcmServerModel != nil {
		return fcmServerModel.Key
	}
	return ""
}

func (d *GormDatabase) UpsertFcmServer(body *model.FcmServer) (*model.FcmServer, error) {
	body.Key = security.Encrypt(body.Key)
	fcmServerModel, _ := d.GetFcmServer()
	if fcmServerModel != nil {
		if err := d.DB.Model(&fcmServerModel).Updates(body).Error; err != nil {
			return nil, err
		}
		return fcmServerModel, nil
	}
	body.UUID = nuuid.ShortUUID("fcm")
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}
