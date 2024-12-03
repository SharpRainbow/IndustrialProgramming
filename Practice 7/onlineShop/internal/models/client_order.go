package models

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"onlineShop/internal/db"
	"time"
)

type ClientOrder struct {
	Id       int              `json:"id"`
	Date     time.Time        `json:"date"`
	Receipt  Receipt          `json:"receipt"`
	Status   Status           `json:"status"`
	ClientId int              `json:"clientId"`
	Products []ProductInOrder `json:"products"`
}

func orderFromRows(rows *sql.Rows) (*ClientOrder, error) {
	var order ClientOrder
	err := rows.Scan(&order.Id, &order.Date, &order.Receipt.Id, &order.Status.Id, &order.ClientId, &order.Receipt.Name, &order.Status.Name)
	return &order, err
}

func ReadClientOrders(db *db.PostgreDb, clientId int) ([]*ClientOrder, error) {
	rows, err := db.ExecutePreparedQuery("SELECT o.*, rt.name, os.name FROM client_order o LEFT JOIN"+
		" receipt_type rt on rt.id = o.receipt_type LEFT JOIN order_status os on os.id = o.status WHERE client = $1", clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*ClientOrder
	for rows.Next() {
		order, err := orderFromRows(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func ReadClientOrderById(db *db.PostgreDb, clientId int, orderId int) ([]*ClientOrder, error) {
	rows, err := db.ExecutePreparedQuery("SELECT o.*, rt.name, os.name FROM client_order o LEFT JOIN"+
		" receipt_type rt on rt.id = o.receipt_type LEFT JOIN order_status os on os.id = o.status WHERE client = $1 AND o.id = $2", clientId, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*ClientOrder
	for rows.Next() {
		order, err := orderFromRows(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func UpdateOrder(db *db.PostgreDb, oid int, uid int, newData *ClientOrder) (int64, error) {
	query := "UPDATE client_order SET"
	var args []interface{}
	argCounter := 1
	if newData.Status.Id > 0 {
		query += " status = $" + fmt.Sprintf("%d", argCounter)
		args = append(args, newData.Status.Id)
		argCounter++
	}
	if newData.Receipt.Id > 0 {
		if argCounter > 1 {
			query += ","
		}
		query += " receipt_type = $" + fmt.Sprintf("%d", argCounter)
		args = append(args, newData.Receipt.Id)
		argCounter++
	}
	query += " WHERE id = $" + fmt.Sprintf("%d", argCounter)
	args = append(args, oid)
	argCounter++
	query += " AND client = $" + fmt.Sprintf("%d", argCounter)
	args = append(args, uid)
	return db.ExecuteInsert(query, args...)
}

func InsertOrder(db *db.PostgreDb, uid int, rid int, products []int, counts []int) error {
	err := db.ExecuteProc("CALL create_order($1, $2, $3, $4)", uid, rid, pq.Array(products), pq.Array(counts))
	return err
}
