package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"onlineShop/internal/db"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func userFromRows(rows *sql.Rows) (*User, error) {
	var user User
	err := rows.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Password)
	return &user, err
}

func ReadUser(db *db.PostgreDb, email string) (*User, error) {
	row, err := db.ExecutePreparedQuery("SELECT * FROM client WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		u := User{}
		err := row.Scan(&u.Id, &u.Name, &u.Phone, &u.Email, &u.Password)
		if err != nil {
			return nil, err
		}
		return &u, nil
	} else {
		return nil, errors.New("unable to find user data")
	}
}

func UpdateUser(db *db.PostgreDb, newData *User, uid int) (int64, error) {
	query := "UPDATE client SET"
	var args []interface{}
	argCounter := 1
	if len(newData.Name) > 0 {
		query += " name = $" + fmt.Sprintf("%d", argCounter)
		args = append(args, newData.Name)
		argCounter++
	}
	if len(newData.Phone) > 0 {
		if argCounter > 1 {
			query += ","
		}
		query += " phone = $" + fmt.Sprintf("%d", argCounter)
		args = append(args, newData.Phone)
		argCounter++
	}
	if len(newData.Password) > 0 {
		if argCounter > 1 {
			query += ","
		}
		query += " password = $" + fmt.Sprintf("%d", argCounter)
		args = append(args, newData.Password)
		argCounter++
	}
	query += " WHERE id = $" + fmt.Sprintf("%d", argCounter)
	args = append(args, uid)
	log.Print(query)
	return db.ExecuteInsert(query, args...)
}
