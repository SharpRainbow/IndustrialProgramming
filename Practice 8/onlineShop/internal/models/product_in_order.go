package models

import (
	"gorm.io/gorm"
)

type ProductInOrder struct {
	ProductId int     `json:"-" gorm:"column:product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductId"`
	Count     int     `json:"count"`
	Price     float64 `json:"price"`
}

func (ProductInOrder) TableName() string {
	return "product_in_order"
}

func ReadProductsFromOrder(db *gorm.DB, userId int, orderId int) (*[]ProductInOrder, error) {
	var pioList []ProductInOrder
	err := db.
		Preload("Product").
		Joins("LEFT JOIN client_order co on product_in_order.order_id = co.id").
		Find(&pioList, "co.id = ? AND co.client = ?", orderId, userId).Error
	return &pioList, err
}

func UpdateProductInOrder(db *gorm.DB, count int, pid int, oid int, uid int) (int64, error) {
	result := db.Model(&ProductInOrder{}).Where("product_id = ? AND order_id = ? "+
		"AND order_id IN (SELECT id FROM client_order WHERE client = ?)", pid, oid, uid).Update("count", count)
	return result.RowsAffected, result.Error
}
