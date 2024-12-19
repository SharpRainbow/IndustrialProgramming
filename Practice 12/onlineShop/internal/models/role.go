package models

import "gorm.io/gorm"

type Role struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (Role) TableName() string {
	return "role_classifier"
}

func ReadRole(db *gorm.DB) (*[]Role, error) {
	var role []Role
	err := db.Find(&role).Error
	return &role, err
}

func ReadRoleById(db *gorm.DB, id int) (*Role, error) {
	var role Role
	err := db.First(&role, id).Error
	return &role, err
}
