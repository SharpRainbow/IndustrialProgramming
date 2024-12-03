package models

import (
	"database/sql"
	"onlineShop/internal/db"
)

type Product struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Description  sql.NullString `json:"description"`
	PictureUrl   sql.NullString `json:"pictureUrl"`
	Category     Category       `json:"category"`
	Manufacturer Manufacturer   `json:"manufacturer"`
	CurrentPrice float64        `json:"price"`
	Parameters   []Parameter    `json:"parameters"`
}

func productFromRows(rows *sql.Rows) (*Product, error) {
	var p Product
	err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.PictureUrl, &p.Category.Id, &p.Manufacturer.Id, &p.CurrentPrice,
		&p.Category.Name, &p.Manufacturer.Id, &p.Manufacturer.Name, &p.Manufacturer.Country.Id,
		&p.Manufacturer.Description, &p.Manufacturer.Country.Name)
	return &p, err
}

func ReadProducts(db *db.PostgreDb) ([]*Product, error) {
	rows, err := db.ExecuteQuery("SELECT p.*, c.name, m.*, c2.name FROM product p LEFT JOIN category c on p.category = c.id" +
		" LEFT JOIN manufacturer m on p.manufacturer = m.id LEFT JOIN country c2 on m.country = c2.id ORDER BY p.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*Product
	for rows.Next() {
		p, err := productFromRows(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func ReadProductById(db *db.PostgreDb, id int) ([]*Product, error) {
	rows, err := db.ExecutePreparedQuery("SELECT p.*, c.name, m.*, c2.name FROM product p LEFT JOIN category c on p.category = c.id"+
		" LEFT JOIN manufacturer m on p.manufacturer = m.id LEFT JOIN country c2 on m.country = c2.id WHERE p.id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*Product
	for rows.Next() {
		p, err := productFromRows(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
