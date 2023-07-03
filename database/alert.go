package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/services/alerts"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"github.com/NubeIO/rubix-os/utils/ttime"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (d *GormDatabase) GetAlert(uuid string) (*model.Alert, error) {
	alertModel := new(model.Alert)
	if err := d.DB.Where("uuid = ? ", uuid).First(&alertModel).Error; err != nil {
		log.Errorf("GetAlert error: %v", err)
		return nil, err
	}
	return alertModel, nil
}

func (d *GormDatabase) GetAlerts() ([]*model.Alert, error) {
	var alertsModel []*model.Alert
	if err := d.DB.Find(&alertsModel).Error; err != nil {
		return nil, err
	}
	return alertsModel, nil
}

func (d *GormDatabase) GetAlertsByHost(hostUUID string) ([]*model.Alert, error) {
	var alertModel []*model.Alert
	alert, err := d.GetAlerts()
	for _, a := range alert {
		if a.HostUUID == hostUUID {
			alertModel = append(alertModel, a)
		}
	}
	return alertModel, err
}

// GetAlertByField returns the object for the given field ie name or nil.
func (d *GormDatabase) GetAlertByField(field string, value string) (*model.Alert, error) {
	var alertModel *model.Alert
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).First(&alertModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return alertModel, nil
}

func (d *GormDatabase) CreateAlert(body *model.Alert) (*model.Alert, error) {
	var err error
	hostUUID := body.HostUUID
	if hostUUID == "" {
		host, err := d.GetFirstHost(api.Args{})
		if err != nil {
			return nil, err
		}
		if host != nil {
			return nil, errors.New(fmt.Sprintf("no host uuid was provided, try uuid: %s, name: %s", host.UUID, host.Name))
		}
		return nil, errors.New(" no host has been added, please add one")
	}
	host, err := d.GetHost(hostUUID, api.Args{})
	if host == nil {
		return nil, errors.New(fmt.Sprintf("host with uuid:%s was not found", hostUUID))
	}
	if body.Status == "" {
		body.Status = string(alerts.Active)
	} else {
		if err = alerts.CheckStatus(body.Status); err != nil {
			return nil, err
		}
	}
	if err = alerts.CheckAlertType(body.Type); err != nil {
		return nil, err
	}
	if err = alerts.CheckSeverity(body.Severity); err != nil {
		return nil, err
	}
	if err = alerts.CheckEntityType(body.EntityType); err != nil {
		return nil, err
	}
	if body.Message == "" {
		body.Message = alerts.AlertTypeMessage(body.Message)
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Alert)
	t := ttime.Now()
	body.CreatedAt = &t
	body.LastUpdated = &t
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateAlertStatus(uuid string, status string) (alert *model.Alert, err error) {
	var query *gorm.DB
	if err = alerts.CheckStatus(status); err != nil {
		return nil, err
	}
	if alerts.CheckStatusClosed(status) { // Move alert to alertClosed table
		alertModel := model.Alert{}
		query = d.DB.Where("uuid = ?", uuid).First(&alertModel)
		if query.Error != nil {
			return nil, query.Error
		}
		ac := model.AlertClosed{
			Alert: alertModel,
		}
		ac.Status = status
		query = d.DB.Create(&ac)
		if query.Error != nil {
			return nil, query.Error
		}
		query = d.DB.Delete(&alertModel)
		if query.Error != nil {
			return nil, query.Error
		}
	} else { // else update alert status
		alert = &model.Alert{}
		query = d.DB.Model(&alert).Where("uuid = ?", uuid).Update("status", status)
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			// look for alert in closed table and re-open
			ac := model.AlertClosed{}
			query = d.DB.Where("uuid = ?", uuid).Find(&ac)
			if query.Error != nil {
				return nil, query.Error
			}
			alert = &ac.Alert
			alert.Status = status
			query = d.DB.Create(&alert)
			if query.Error != nil {
				return nil, query.Error
			}
			query = d.DB.Delete(&ac)
			if query.Error != nil {
				return nil, query.Error
			}
		}
	}
	t := ttime.Now()
	alert.LastUpdated = &t
	return alert, query.Error
}

func (d *GormDatabase) DeleteAlert(uuid string) (*interfaces.Message, error) {
	alertModel := new(model.Alert)
	query := d.DB.Where("uuid = ? ", uuid).Delete(&alertModel)
	return d.deleteResponse(query)
}

func (d *GormDatabase) DropAlerts() (*interfaces.Message, error) {
	var alertModel *model.Alert
	query := d.DB.Where("1 = 1").Delete(&alertModel)
	return d.deleteResponse(query)
}
