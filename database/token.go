package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type Token struct {
	*model.Token
}

func (d *GormDatabase) GetTokens() ([]*model.Token, error) {
	var tokens []*model.Token
	if err := d.DB.Omit("Token").Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

// CreateToken create a thing
func (d *GormDatabase) CreateToken(body *model.Token) (*model.Token, error) {
	if err := d.DB.Create(body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

// UpdateToken update a thing
func (d *GormDatabase) UpdateToken(body *model.Token) (*model.Token, error) {
	var tokenModel *model.Token
	if err := d.DB.Where("name = ?", body.Name).First(&tokenModel).Error; err != nil {
		return nil, err
	}
	if err := d.DB.Omit("Token").Model(&tokenModel).Updates(body).Error; err != nil {
		return nil, err
	}
	tokenModel.Token = ""
	return tokenModel, nil
}
