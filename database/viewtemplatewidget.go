package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"sync"
)

func (d *GormDatabase) GetViewTemplateWidgets(args argspkg.Args) ([]*model.ViewTemplateWidget, error) {
	var viewTemplateWidgetsModel []*model.ViewTemplateWidget
	query := d.buildViewTemplateWidgetQuery(args)
	if err := query.Find(&viewTemplateWidgetsModel).Error; err != nil {
		return nil, err
	}
	return viewTemplateWidgetsModel, nil
}

func (d *GormDatabase) GetViewTemplateWidget(uuid string, args argspkg.Args) (*model.ViewTemplateWidget, error) {
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
	syncAfterUpdate := false
	if body.ThingType == "" {
		body.ThingType = viewTemplateWidgetModel.ThingType
	}
	switch body.ThingType {
	case model.ThingType.Point:
		syncAfterUpdate = body.NetworkName != viewTemplateWidgetModel.NetworkName ||
			body.DeviceName != viewTemplateWidgetModel.DeviceName || body.ThingName != viewTemplateWidgetModel.ThingName
	case model.ThingType.Schedule:
		syncAfterUpdate = body.ThingName != viewTemplateWidgetModel.ThingName
	}
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
	viewTemplateWidget, err := d.GetViewTemplateWidget(uuid, argspkg.Args{WithViewTemplateWidgetPointers: true})
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
			switch vtWidget.ThingType {
			case model.ThingType.Point:
				point, connectionErr, requestErr := cli.GetPointByName(*vtWidget.NetworkName, *vtWidget.DeviceName,
					vtWidget.ThingName)
				if connectionErr != nil || requestErr != nil {
					return
				}
				vtwPointer.DeviceUUID = nstring.New(point.DeviceUUID)
				vtwPointer.ThingUUID = point.UUID
				_, _ = d.UpdateViewTemplateWidgetPointer(vtwPointer.UUID, vtwPointer)
			case model.ThingType.Schedule:
				schedule, connectionErr, requestErr := cli.GetScheduleByNameV2(vtWidget.ThingName)
				if connectionErr != nil || requestErr != nil {
					return
				}
				vtwPointer.DeviceUUID = nil
				vtwPointer.ThingUUID = schedule.UUID
				_, _ = d.UpdateViewTemplateWidgetPointer(vtwPointer.UUID, vtwPointer)
			}
		}(viewTemplateWidget, viewTemplateWidgetPointer)
	}
	wg.Wait()
}
