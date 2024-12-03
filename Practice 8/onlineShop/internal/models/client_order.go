package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type ClientOrder struct {
	Id        int              `json:"id"`
	Date      time.Time        `json:"date"`
	ReceiptId int              `json:"-" gorm:"column:receipt_type"`
	Receipt   Receipt          `json:"receipt" gorm:"foreignKey:ReceiptId"`
	StatusId  int              `json:"-" gorm:"column:status"`
	Status    Status           `json:"status" gorm:"foreignKey:StatusId"`
	ClientId  int              `json:"clientId" gorm:"column:client"`
	Products  []ProductInOrder `json:"products" gorm:"-"`
}

func (ClientOrder) TableName() string {
	return "client_order"
}

func ReadClientOrders(db *gorm.DB, clientId int) (*[]ClientOrder, error) {
	var orders []ClientOrder
	err := db.Preload("Status").Preload("Receipt").Find(&orders, "client = ?", clientId).Error
	return &orders, err
}

func ReadClientOrderById(db *gorm.DB, clientId int, orderId int) (*[]ClientOrder, error) {
	var orders []ClientOrder
	err := db.Preload("Status").Preload("Receipt").Find(&orders, "client = ? AND id = ?", clientId, orderId).Error
	return &orders, err
}

func UpdateOrder(db *gorm.DB, oid int, uid int, newData *ClientOrder) (int64, error) {
	updates := map[string]interface{}{}
	if newData.Status.Id > 0 {
		updates["status"] = newData.Status.Id
	}
	if newData.Receipt.Id > 0 {
		updates["receipt_type"] = newData.Receipt.Id
	}
	result := db.Model(&ClientOrder{}).
		Where("id = ? AND client = ?", oid, uid).
		Updates(updates)
	return result.RowsAffected, result.Error
}

func InsertOrder(db *gorm.DB, uid int, rid int, products []int, counts []int) error {
	err := db.Exec("CALL create_order($1, $2, $3, $4)", uid, rid, pq.Array(products), pq.Array(counts)).Error
	return err
}
