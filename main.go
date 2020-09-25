package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const connectionTemplate = "%s:%s@(%s:%s)/%s?parseTime=true&tls=%t&multiStatements=true"

func connect() (*sql.DB, error) {
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

func main() {
	if err := transaction(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type user struct {
	id        int
	firstName string
	lastName  string
}

func (u *user) String() string {
	return fmt.Sprintf("id: %d, first_name: %s, last_name: %s", u.id, u.firstName, u.lastName)
}

func exec(db DB, statement string) error {
	result, err := db.Exec(statement)
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func selectUsersRecords(db DB, query string) ([]*user, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*user, 0)

	for rows.Next() {
		u := &user{}
		if err := rows.Scan(&u.id, &u.firstName, &u.lastName); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if !rows.NextResultSet() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return users, nil
}
