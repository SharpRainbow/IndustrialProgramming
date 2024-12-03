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

func ReadCountries(db *gorm.DB) (*[]Country, error) {
	var countries []Country
	err := db.Find(&countries).Error
	return &countries, err
}
