package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetViewWidget(uuid string) (*model.ViewWidget, error) {
	var viewWidgetModel *model.ViewWidget
	if err := d.DB.Where("uuid = ?", uuid).First(&viewWidgetModel).Error; err != nil {
		return nil, err
	}
	return viewWidgetModel, nil
}

func (d *GormDatabase) GetViewWidgetByViewUUID(viewUUID string) ([]*model.ViewWidget, error) {
	var viewWidgetsModel []*model.ViewWidget
	if err := d.DB.Where("view_uuid = ?", viewUUID).Find(&viewWidgetsModel).Error; err != nil {
		return nil, err
	}
	return viewWidgetsModel, nil
}

func (d *GormDatabase) CreateViewWidget(body *model.ViewWidget) (*model.ViewWidget, error) {
	if err := d.mapBeforeCreateUpdateViewWidget(body); err != nil {
		return nil, err
	}
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.ViewWidget)
	body.Config = marshalJson(body.Config)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateViewWidget(uuid string, body *model.ViewWidget) (*model.ViewWidget, error) {
	viewWidgetModel, err := d.GetViewWidget(uuid)
	if err != nil {
		return nil, err
	}
	if body.HostUUID != body.HostUUID || viewWidgetModel.PointUUID != viewWidgetModel.PointUUID {
		if err := d.mapBeforeCreateUpdateViewWidget(body); err != nil {
			return nil, err
		}
	}
	if body.Config != nil {
		body.Config = marshalJson(body.Config)
	}
	d.DB.Model(&viewWidgetModel).Updates(&body)
	return viewWidgetModel, nil
}

func (d *GormDatabase) DeleteViewWidget(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.ViewWidget{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteViewWidgetsByViewUUID(viewUUID string) (bool, error) {
	query := d.DB.Where("view_uuid = ?", viewUUID).Delete(&model.ViewWidget{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) mapBeforeCreateUpdateViewWidget(body *model.ViewWidget) error {
	host, err := d.ResolveHost(body.HostUUID, "")
	if err != nil {
		return err
	}
	cli := client.NewClient(host.IP, host.Port, host.ExternalToken)
	point, err := cli.GetPointWithParent(body.PointUUID)
	if err != nil {
		return err
	}
	body.NetworkName = point.NetworkName
	body.DeviceName = point.DeviceName
	body.PointName = point.Name
	body.DeviceUUID = point.DeviceUUID
	return nil
}
