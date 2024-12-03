package models

import (
	"gorm.io/gorm"
)

type Product struct {
	Id             int          `json:"id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	PictureUrl     string       `json:"pictureUrl" gorm:"column:picture_link"`
	CategoryId     int          `json:"-" gorm:"column:category"`
	Category       Category     `json:"category" gorm:"foreignKey:CategoryId"`
	ManufacturerId int          `json:"-" gorm:"column:manufacturer"`
	Manufacturer   Manufacturer `json:"manufacturer" gorm:"foreignKey:ManufacturerId"`
	CurrentPrice   float64      `json:"price" gorm:"column:current_price"`
}

func (Product) TableName() string {
	return "product"
}

func ReadProducts(db *gorm.DB) (*[]Product, error) {
	var products []Product
	err := db.Preload("Manufacturer").Preload("Category").Find(&products).Error
	return &products, err
}

func ReadProductById(db *gorm.DB, id int) (*[]Product, error) {
	var products []Product
	err := db.Preload("Manufacturer").Preload("Category").Find(&products, id).Error
	return &products, err
}
