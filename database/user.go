package database

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type User struct {
	*model.User
}

// GetUser get it
func (d *GormDatabase) GetUser() (*model.User, error) {
	var userModel *model.User
	query := d.DB.First(&userModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return userModel, nil
}

// UpdateUser update or create if not found a thing
func (d *GormDatabase) UpdateUser(body *model.User) (*model.User, error) {
	var userModel *model.User
	query := d.DB.First(&userModel)
	if userModel.Username == "" {
		if err := d.DB.Create(&body).Error; err != nil {
			return nil, err
		}
		return body, nil
	}
	query = d.DB.Model(&userModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return userModel, nil
}
