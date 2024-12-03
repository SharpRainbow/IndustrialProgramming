package models

import (
	"database/sql"
	"onlineShop/internal/db"
)

type Status struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func statusFromRows(rows *sql.Rows) (*Status, error) {
	var status Status
	err := rows.Scan(&status.Id, &status.Name)
	return &status, err
}

func ReadStatuses(db *db.PostgreDb) ([]*Status, error) {
	rows, err := db.ExecuteQuery("SELECT * FROM order_status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var statuses []*Status
	for rows.Next() {
		status, err := statusFromRows(rows)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}
