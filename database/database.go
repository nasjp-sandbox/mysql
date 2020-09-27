package database

import (
	"database/sql"
	"fmt"

	// for sql driver.
	_ "github.com/go-sql-driver/mysql"
)

const connectionTemplate = "%s:%s@(%s:%s)/%s?parseTime=true&tls=%t&multiStatements=true"

type database interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		connectionTemplate,
		"root",
		"",
		"db",
		"3306",
		"",
		false,
	))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS app"); err != nil {
		return nil, err
	}

	return db, nil
}

func Exec(db database, statement string, args ...interface{}) error {
	result, err := db.Exec(statement, args...)
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	return nil
}
