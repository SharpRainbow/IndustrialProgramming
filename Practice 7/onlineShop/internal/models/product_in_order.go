package models

import (
	"database/sql"
	"onlineShop/internal/db"
)

type ProductInOrder struct {
	Product Product `json:"product"`
	Count   int     `json:"count"`
	Price   float64 `json:"price"`
}

func pioFromRows(rows *sql.Rows) (*ProductInOrder, error) {
	var p ProductInOrder
	err := rows.Scan(&p.Product.Id, &p.Product.Name, &p.Product.Description, &p.Product.PictureUrl,
		&p.Product.Category.Id, &p.Product.Manufacturer.Id, &p.Product.CurrentPrice, &p.Count, &p.Price)
	return &p, err
}

func ReadProductsFromOrder(db *db.PostgreDb, userId int, orderId int) ([]*ProductInOrder, error) {
	rows, err := db.ExecutePreparedQuery("SELECT p.*, pio.count, pio.price FROM product_in_order pio"+
		" LEFT JOIN product p on p.id = pio.product_id LEFT JOIN client_order co on pio.order_id = co.id WHERE order_id = $1 AND client = $2", orderId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pioList []*ProductInOrder
	for rows.Next() {
		pio, err := pioFromRows(rows)
		if err != nil {
			return nil, err
		}
		pioList = append(pioList, pio)
	}
	return pioList, nil
}

func UpdateProductInOrder(db *db.PostgreDb, count int, pid int, oid int, uid int) (int64, error) {
	return db.ExecuteInsert("UPDATE product_in_order SET count = $1 WHERE product_id = $2"+
		" AND order_id = $3 AND order_id IN (SELECT id FROM client_order WHERE client = $4)", count,
		pid, oid, uid)
}
