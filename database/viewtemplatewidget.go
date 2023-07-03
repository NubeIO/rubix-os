package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"sync"
)

func (d *GormDatabase) GetViewTemplateWidgets(args api.Args) ([]*model.ViewTemplateWidget, error) {
	var viewTemplateWidgetsModel []*model.ViewTemplateWidget
	query := d.buildViewTemplateWidgetQuery(args)
	if err := query.Find(&viewTemplateWidgetsModel).Error; err != nil {
		return nil, err
	}
	return viewTemplateWidgetsModel, nil
}

func (d *GormDatabase) GetViewTemplateWidget(uuid string, args api.Args) (*model.ViewTemplateWidget, error) {
	var viewTemplateWidgetModel *model.ViewTemplateWidget
	query := d.buildViewTemplateWidgetQuery(args)
	if err := query.Where("uuid = ?", uuid).First(&viewTemplateWidgetModel).Error; err != nil {
		return nil, err
	}
	return viewTemplateWidgetModel, nil
}

func (d *GormDatabase) CreateViewTemplateWidget(body *model.ViewTemplateWidget) (*model.ViewTemplateWidget, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.ViewTemplateWidget)
	body.Config = marshalJson(body.Config)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateViewTemplateWidget(uuid string, body *model.ViewTemplateWidget) (
	*model.ViewTemplateWidget, error) {
	var viewTemplateWidgetModel *model.ViewTemplateWidget
	if err := d.DB.Where("uuid = ?", uuid).First(&viewTemplateWidgetModel).Error; err != nil {
		return nil, err
	}
	syncAfterUpdate := body.NetworkName != viewTemplateWidgetModel.NetworkName ||
		body.DeviceName != viewTemplateWidgetModel.DeviceName || body.PointName != viewTemplateWidgetModel.PointName
	if body.Config != nil {
		body.Config = marshalJson(body.Config)
	}
	if err := d.DB.Model(&viewTemplateWidgetModel).Updates(body).Error; err != nil {
		return nil, err
	}
	if syncAfterUpdate {
		d.syncAfterUpdateViewTemplateWidget(uuid)
	}
	return viewTemplateWidgetModel, nil
}

func (d *GormDatabase) DeleteViewTemplateWidget(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.ViewTemplateWidget{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) syncAfterUpdateViewTemplateWidget(uuid string) {
	viewTemplateWidget, err := d.GetViewTemplateWidget(uuid, api.Args{WithViewTemplateWidgetPointers: true})
	if err != nil {
		return
	}
	wg := &sync.WaitGroup{}
	for _, viewTemplateWidgetPointer := range viewTemplateWidget.ViewTemplateWidgetPointers {
		wg.Add(1)
		go func(vtWidget *model.ViewTemplateWidget, vtwPointer *model.ViewTemplateWidgetPointer) {
			defer wg.Done()
			host, err := d.ResolveHost(vtwPointer.HostUUID, "")
			if err != nil {
				return
			}
			cli := client.NewClient(host.IP, host.Port, host.ExternalToken)
			point, err, err1 := cli.GetPointByName(vtWidget.NetworkName, vtWidget.DeviceName, vtWidget.PointName)
			if err != nil || err1 != nil {
				return
			}
			vtwPointer.DeviceUUID = point.DeviceUUID
			vtwPointer.PointUUID = point.UUID
			_, _ = d.UpdateViewTemplateWidgetPointer(vtwPointer.UUID, vtwPointer)
		}(viewTemplateWidget, viewTemplateWidgetPointer)
	}
	wg.Wait()
}
