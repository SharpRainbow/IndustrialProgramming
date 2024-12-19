package models

import (
	"gorm.io/gorm"
)

type Product struct {
	Id             int           `json:"id"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	PictureUrl     string        `json:"pictureUrl" gorm:"column:picture_link"`
	CategoryId     int           `json:"-" gorm:"column:category"`
	Category       *Category     `json:"category,omitempty" gorm:"foreignKey:CategoryId"`
	ManufacturerId int           `json:"-" gorm:"column:manufacturer"`
	Manufacturer   *Manufacturer `json:"manufacturer,omitempty" gorm:"foreignKey:ManufacturerId"`
	CurrentPrice   float64       `json:"price" gorm:"column:current_price"`
}

var pSortFields = []string{"id", "name", "current_price", "category"}

func (Product) TableName() string {
	return "product"
}

func ProductSortFieldValid(sort string) bool {
	for _, field := range pSortFields {
		if field == sort {
			return true
		}
	}
	return false
}

func ReadProducts(db *gorm.DB, limit int, offset int, name string, sort string) (*[]Product, error) {
	var products []Product
	query := db.Preload("Manufacturer").Preload("Category").Limit(limit).Offset(offset)
	if len(name) > 0 {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	err := query.Order(sort).Find(&products).Error
	return &products, err
}

func ReadProductById(db *gorm.DB, id int) (*[]Product, error) {
	var products []Product
	err := db.Preload("Manufacturer").Preload("Category").Find(&products, id).Error
	return &products, err
}
