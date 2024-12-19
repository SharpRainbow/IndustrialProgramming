package models

import (
	"gorm.io/gorm"
)

type Manufacturer struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	CountryId   int      `json:"-" gorm:"column:country"`
	Country     *Country `json:"country,omitempty" gorm:"foreignKey:CountryId"`
	Description string   `json:"description"`
}

func (Manufacturer) TableName() string {
	return "manufacturer"
}

var mSortFields = []string{"id", "name"}

func ManufacturerSortFieldValid(sort string) bool {
	for _, field := range mSortFields {
		if field == sort {
			return true
		}
	}
	return false
}

func ReadManufacturers(db *gorm.DB, limit int, offset int, name string, sort string) (*[]Manufacturer, error) {
	var manufacturers []Manufacturer
	query := db.Preload("Country").Limit(limit).Offset(offset)
	if len(name) > 0 {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	err := query.Order(sort).Find(&manufacturers).Error
	return &manufacturers, err
}

func ReadManufacturerById(db *gorm.DB, id int) (*[]Manufacturer, error) {
	var manufacturers []Manufacturer
	err := db.Preload("Country").Find(&manufacturers, id).Error
	return &manufacturers, err
}
