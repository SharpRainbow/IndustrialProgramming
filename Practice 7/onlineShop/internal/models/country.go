package models

import (
	"database/sql"
	"onlineShop/internal/db"
)

type Country struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func countryFromRows(rows *sql.Rows) (*Country, error) {
	var country Country
	err := rows.Scan(&country.Id, &country.Name)
	return &country, err
}

func ReadCountries(db *db.PostgreDb) ([]*Country, error) {
	rows, err := db.ExecuteQuery("SELECT * FROM country")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var countries []*Country
	for rows.Next() {
		category, err := countryFromRows(rows)
		if err != nil {
			return nil, err
		}
		countries = append(countries, category)
	}
	return countries, nil
}
