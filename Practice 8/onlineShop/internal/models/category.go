package models

import (
	"gorm.io/gorm"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (Category) TableName() string {
	return "category"
}

func ReadCategories(db *gorm.DB) (*[]Category, error) {
	var categories []Category
	err := db.Find(&categories).Error
	return &categories, err
}

func ReadCategoryById(db *gorm.DB, id int) (*[]Category, error) {
	var categories *[]Category
	err := db.Find(&categories, id).Error
	return categories, err
}

func ReadCategoryOfManufacturer(db *gorm.DB, manufacturerId int) (*[]Category, error) {
	var categories []Category
	err := db.Joins("LEFT JOIN category_of_manufacturer com ON category.id = com.category_id").
		Where("com.manufacturer_id = ?", manufacturerId).
		Find(&categories).Error
	return &categories, err
}
