package models

import (
	"gorm.io/gorm"
)

type Receipt struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (Receipt) TableName() string {
	return "receipt_type"
}

func ReadReceiptTypes(db *gorm.DB) (*[]Receipt, error) {
	var receiptTypes []Receipt
	err := db.Find(&receiptTypes).Error
	return &receiptTypes, err
}
