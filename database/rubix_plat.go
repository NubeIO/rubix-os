package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


// GetRubixPlat returns all networks.
func (d *GormDatabase) GetRubixPlat() (*model.RubixPlat, error) {
	var rubixPlatModel *model.RubixPlat
		query := d.DB.Find(&rubixPlatModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return rubixPlatModel, nil
}


// CreateRubixPlat creates a device.
func (d *GormDatabase) CreateRubixPlat(body *model.RubixPlat) (*model.RubixPlat, error) {
	body.GlobalUuid = utils.MakeTopicUUID(model.CommonNaming.Network)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}



// UpdateRubixPlat update it.
func (d *GormDatabase) UpdateRubixPlat(body *model.RubixPlat) (*model.RubixPlat, error) {
	var rubixPlatModel *model.RubixPlat
	query := d.DB.Where("id >= ?", 0).Find(&rubixPlatModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&rubixPlatModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return rubixPlatModel, nil

}




// DeleteRubixPlat delete a network.
func (d *GormDatabase) DeleteRubixPlat() (bool, error) {
	var rubixPlatModel *model.RubixPlat
	query := d.DB.Delete(&rubixPlatModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}
