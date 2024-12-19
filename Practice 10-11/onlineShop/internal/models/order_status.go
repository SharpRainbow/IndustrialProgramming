package models

import (
	"gorm.io/gorm"
)

type Status struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (Status) TableName() string {
	return "order_status"
}

func ReadStatuses(db *gorm.DB) (*[]Status, error) {
	var statuses []Status
	err := db.Find(&statuses).Error
	return &statuses, err
}
