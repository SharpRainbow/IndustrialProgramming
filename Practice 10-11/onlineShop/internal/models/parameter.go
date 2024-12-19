package models

import (
	"gorm.io/gorm"
)

type Parameter struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

func (Parameter) TableName() string {
	return "parameter"
}

func ReadParametersOfCategory(db *gorm.DB, categoryId int) (*[]Parameter, error) {
	var params []Parameter
	err := db.
		Select("parameter.id, parameter.name, parameter.description").
		Joins("LEFT JOIN parameter_of_category poc on parameter.id = poc.parameter_id").
		Find(&params, "poc.category_id = ?", categoryId).Error
	return &params, err
}

func ReadParametersOfProduct(db *gorm.DB, productId int) (*[]Parameter, error) {
	var params []Parameter
	err := db.
		Select("parameter.id, parameter.name, parameter.description, pp.value").
		Joins("RIGHT JOIN parameter_of_product pp on parameter.id = pp.parameter_id").
		Find(&params, "pp.product_id = ?", productId).Error
	return &params, err
}
