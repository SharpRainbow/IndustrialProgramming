package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgreDb struct {
	db *sql.DB
}

func (conn PostgreDb) ExecuteInsert(query string, params ...any) (int64, error) {
	rows, err := conn.db.Exec(query, params...)
	if err != nil {
		return 0, err
	}
	return rows.RowsAffected()
}

func (conn PostgreDb) ExecuteProc(query string, params ...any) error {
	_, err := conn.db.Exec(query, params...)
	if err != nil {
		return err
	}
	return nil
}

func (conn PostgreDb) ExecuteQuery(query string) (*sql.Rows, error) {
	rows, err := conn.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (conn PostgreDb) ExecutePreparedQuery(query string, params ...any) (*sql.Rows, error) {
	stmnt, err := conn.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmnt.Close()
	rows, err := stmnt.Query(params...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

var connection = &PostgreDb{}

func GetDB(cfg *Config) (*PostgreDb, error) {
	if connection.db != nil && connection.db.Ping() == nil {
		return connection, nil
	}
	var err error
	connection.db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbName, cfg.DbPassword))
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func CloseDB() {
	if connection != nil && connection.db != nil {
		connection.db.Close()
	}
}
