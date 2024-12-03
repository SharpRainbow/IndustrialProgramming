package models

import (
	"database/sql"
	"onlineShop/internal/db"
)

type Parameter struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

func parameterOfCategoryFromRows(rows *sql.Rows) (*Parameter, error) {
	var param Parameter
	err := rows.Scan(&param.Id, &param.Name, &param.Description)
	return &param, err
}

func parameterOfProductFromRows(rows *sql.Rows) (*Parameter, error) {
	var param Parameter
	err := rows.Scan(&param.Id, &param.Name, &param.Description, &param.Value)
	return &param, err
}

func ReadParametersOfCategory(db *db.PostgreDb, categoryId int) ([]*Parameter, error) {
	rows, err := db.ExecutePreparedQuery("SELECT p.* FROM parameter p LEFT JOIN parameter_of_category poc on p.id = poc.parameter_id WHERE category_id = $1", categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var params []*Parameter
	for rows.Next() {
		param, err := parameterOfCategoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return params, nil
}

func ReadParametersOfProduct(db *db.PostgreDb, productId int) ([]*Parameter, error) {
	rows, err := db.ExecutePreparedQuery("SELECT p.*, pp.value FROM parameter_of_product pp LEFT JOIN parameter p on p.id = pp.parameter_id WHERE product_id = $1", productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var params []*Parameter
	for rows.Next() {
		param, err := parameterOfProductFromRows(rows)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return params, nil
}
