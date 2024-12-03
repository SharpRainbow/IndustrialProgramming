package models

import (
	"database/sql"
	"onlineShop/internal/db"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func categoryFromRows(rows *sql.Rows) (*Category, error) {
	var category Category
	err := rows.Scan(&category.Id, &category.Name)
	return &category, err
}

func ReadCategories(db *db.PostgreDb) ([]*Category, error) {
	rows, err := db.ExecuteQuery("SELECT * FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*Category
	for rows.Next() {
		category, err := categoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func ReadCategoryById(db *db.PostgreDb, id int) ([]*Category, error) {
	rows, err := db.ExecutePreparedQuery("SELECT * FROM category WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*Category
	for rows.Next() {
		category, err := categoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func ReadCategoryOfManufacturer(db *db.PostgreDb, manufacturerId int) ([]*Category, error) {
	rows, err := db.ExecutePreparedQuery("SELECT c.* FROM category c LEFT JOIN category_of_manufacturer com on c.id = com.category_id WHERE manufacturer_id = $1", manufacturerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*Category
	for rows.Next() {
		category, err := categoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
