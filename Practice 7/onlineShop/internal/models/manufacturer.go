package models

import (
	"database/sql"
	"onlineShop/internal/db"
)

type Manufacturer struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Country     Country `json:"country"`
	Description string  `json:"description"`
}

func manufacturerFromRows(rows *sql.Rows) (*Manufacturer, error) {
	var manufacturer Manufacturer
	err := rows.Scan(&manufacturer.Id, &manufacturer.Name, &manufacturer.Country.Id, &manufacturer.Country.Name, &manufacturer.Description)
	return &manufacturer, err
}

func ReadManufacturers(db *db.PostgreDb) ([]*Manufacturer, error) {
	rows, err := db.ExecuteQuery("SELECT m.*, c.name FROM manufacturer m LEFT JOIN country c ON m.country = c.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var manufacturers []*Manufacturer
	for rows.Next() {
		manufacturer, err := manufacturerFromRows(rows)
		if err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, manufacturer)
	}
	return manufacturers, nil
}

func ReadManufacturerById(db *db.PostgreDb, id int) ([]*Manufacturer, error) {
	rows, err := db.ExecutePreparedQuery("SELECT m.*, c.name FROM manufacturer m LEFT JOIN country c ON m.country = c.id WHERE m.id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var manufacturers []*Manufacturer
	for rows.Next() {
		category, err := manufacturerFromRows(rows)
		if err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, category)
	}
	return manufacturers, rows.Err()
}
