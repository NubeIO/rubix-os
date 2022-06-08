package database

func (d *GormDatabase) ValidateToken(token string) (bool, error) {
	var count int64
	if err := d.DB.Table("tokens").Where("token = ?", token).Where("blocked IS NOT TRUE").
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
