package models

import (
	"database/sql"
	"onlineShop/internal/db"
)

type Receipt struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func receiptFromRows(rows *sql.Rows) (*Receipt, error) {
	var receipt Receipt
	err := rows.Scan(&receipt.Id, &receipt.Name)
	return &receipt, err
}

func ReadReceiptTypes(db *db.PostgreDb) ([]*Receipt, error) {
	rows, err := db.ExecuteQuery("SELECT * FROM receipt_type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var receiptTypes []*Receipt
	for rows.Next() {
		receipt, err := receiptFromRows(rows)
		if err != nil {
			return nil, err
		}
		receiptTypes = append(receiptTypes, receipt)
	}
	return receiptTypes, nil
}
