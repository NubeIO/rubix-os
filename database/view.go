package database

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"sync"
)

func (d *GormDatabase) GetViews() ([]*model.View, error) {
	var viewsModel []*model.View
	query := d.buildViewQuery()
	if err := query.Find(&viewsModel).Error; err != nil {
		return nil, err
	}
	return viewsModel, nil
}

func (d *GormDatabase) GetView(uuid string) (*model.View, error) {
	var viewModel *model.View
	query := d.buildViewQuery()
	if err := query.Where("uuid = ?", uuid).First(&viewModel).Error; err != nil {
		return nil, err
	}
	return viewModel, nil
}

func (d *GormDatabase) GetViewsByUUIDs(uuids []*string) ([]*model.View, error) {
	var viewsModel []*model.View
	query := d.buildViewQuery()
	if err := query.Where("uuid IN ?", uuids).Find(&viewsModel).Error; err != nil {
		return nil, err
	}
	return viewsModel, nil
}

func (d *GormDatabase) GetViewsByMemberUsername(memberUsername string) ([]*model.View, error) {
	var viewsModel []*model.View
	query := d.DB.Distinct("views.*").
		Joins("JOIN team_views ON team_views.view_uuid = views.uuid").
		Joins("JOIN teams ON teams.uuid = team_views.team_uuid").
		Joins("JOIN team_members ON team_members.team_uuid = teams.uuid").
		Joins("JOIN members ON members.uuid = team_members.member_uuid").
		Where("members.username = ?", memberUsername)
	if err := query.Find(&viewsModel).Error; err != nil {
		return nil, err
	}
	return viewsModel, nil
}

func (d *GormDatabase) CreateView(body *model.View) (*model.View, error) {
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.View)
	body.Name = name
	body.WidgetConfig = marshalJson(body.WidgetConfig)
	body.Theme = marshalJson(body.Theme)
	if body.LocationUUID == nil && body.GroupUUID == nil && body.HostUUID == nil {
		return nil, errors.New("view should assign either to the location, group or host")
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateView(uuid string, body *model.View) (*model.View, error) {
	var viewModel *model.View
	if err := d.DB.Where("uuid = ?", uuid).First(&viewModel).Error; err != nil {
		return nil, err
	}
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	body.Name = name
	body.ViewTemplateUUID = viewModel.ViewTemplateUUID
	if body.WidgetConfig != nil {
		body.WidgetConfig = marshalJson(body.WidgetConfig)
	}
	if body.Theme != nil {
		body.Theme = marshalJson(body.Theme)
	}
	if err = d.DB.Model(&viewModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return viewModel, nil
}

func (d *GormDatabase) DeleteView(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.View{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) GenerateViewTemplate(uuid string, name string) (bool, error) {
	viewWidgets, _ := d.GetViewWidgetByViewUUID(uuid)
	if len(viewWidgets) == 0 {
		return false, errors.New("no widgets has been added. please add one. ")
	}

	_, _ = d.DeleteViewTemplateWidgetPointerByViewUUID(uuid)

	viewTemplate, err := d.CreateViewTemplate(&model.ViewTemplate{CommonNameUnique: model.CommonNameUnique{Name: name}})
	if err != nil {
		return false, err
	}
	wg := &sync.WaitGroup{}
	for _, viewWidget := range viewWidgets {
		wg.Add(1)
		go func(vWidget *model.ViewWidget) {
			defer wg.Done()
			viewTemplateWidget := &model.ViewTemplateWidget{
				ViewTemplateUUID: viewTemplate.UUID,
				Name:             vWidget.Name,
				Order:            vWidget.Order,
				X:                vWidget.X,
				Y:                vWidget.Y,
				Type:             vWidget.Type,
				Config:           vWidget.Config,
				NetworkName:      vWidget.NetworkName,
				DeviceName:       vWidget.DeviceName,
				PointName:        vWidget.PointName,
				ViewTemplateWidgetPointers: []*model.ViewTemplateWidgetPointer{
					{
						CommonUUID: model.CommonUUID{UUID: nuuid.MakeTopicUUID(model.CommonNaming.ViewTemplateWidgetPointer)},
						ViewUUID:   uuid,
						HostUUID:   vWidget.HostUUID,
						DeviceUUID: vWidget.DeviceUUID,
						PointUUID:  vWidget.PointUUID,
					},
				},
			}
			_, err = d.CreateViewTemplateWidget(viewTemplateWidget)
		}(viewWidget)
	}
	wg.Wait()
	_, err = d.DeleteViewWidgetsByViewUUID(uuid)
	d.updateViewTemplateUUID(uuid, viewTemplate.UUID)
	return true, nil
}

func (d *GormDatabase) AssignViewTemplate(uuid string, viewTemplateUUID string, hostUUID string) (bool, error) {
	viewTemplate, err := d.GetViewTemplate(viewTemplateUUID)
	if err != nil {
		return false, err
	}
	host, err := d.ResolveHost(hostUUID, "")
	if err != nil {
		return false, err
	}
	cli := client.NewClient(host.IP, host.Port, host.ExternalToken)

	_, _ = d.DeleteViewTemplateWidgetPointerByViewUUID(uuid)
	wg := &sync.WaitGroup{}
	for _, viewTemplateWidget := range viewTemplate.ViewTemplateWidgets {
		wg.Add(1)
		go func(vtWidget *model.ViewTemplateWidget) {
			defer wg.Done()
			point, err, err1 := cli.GetPointByName(vtWidget.NetworkName, vtWidget.DeviceName, vtWidget.PointName)
			if err != nil || err1 != nil {
				return
			}
			viewTemplateWidgetPointer := &model.ViewTemplateWidgetPointer{
				ViewTemplateWidgetUUID: vtWidget.UUID,
				ViewUUID:               uuid,
				HostUUID:               hostUUID,
				DeviceUUID:             point.DeviceUUID,
				PointUUID:              point.UUID,
			}
			_, err = d.CreateViewTemplateWidgetPointer(viewTemplateWidgetPointer)
		}(viewTemplateWidget)
	}
	wg.Wait()
	d.updateViewTemplateUUID(uuid, viewTemplate.UUID)
	return true, nil
}

func (d *GormDatabase) updateViewTemplateUUID(uuid string, templateUUID string) {
	d.DB.Model(&model.View{}).Where("uuid = ?", uuid).Update("view_template_uuid", templateUUID)
}
