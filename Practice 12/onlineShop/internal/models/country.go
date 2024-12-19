package models

import (
	"gorm.io/gorm"
)

type Country struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (Country) TableName() string {
	return "country"
}

func ReadCountries(db *gorm.DB, limit int, offset int) (*[]Country, error) {
	var countries []Country
	err := db.Limit(limit).Offset(offset).Find(&countries).Error
	return &countries, err
}
