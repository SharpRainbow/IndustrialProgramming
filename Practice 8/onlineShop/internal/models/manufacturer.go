package models

import (
	"gorm.io/gorm"
)

type Manufacturer struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	CountryId   int     `json:"-" gorm:"column:country"`
	Country     Country `json:"country" gorm:"foreignKey:CountryId"`
	Description string  `json:"description"`
}

func (Manufacturer) TableName() string {
	return "manufacturer"
}

func ReadManufacturers(db *gorm.DB) (*[]Manufacturer, error) {
	var manufacturers []Manufacturer
	db.Preload("Country").Find(&manufacturers)
	return &manufacturers, nil
}

func ReadManufacturerById(db *gorm.DB, id int) (*[]Manufacturer, error) {
	var manufacturers []Manufacturer
	err := db.Preload("Country").Find(&manufacturers, id).Error
	return &manufacturers, err
}
