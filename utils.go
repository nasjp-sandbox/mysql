package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const connectionTemplate = "%s:%s@(%s:%s)/%s?parseTime=true&tls=%t&multiStatements=true"

type database interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

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

func exec(db database, statement string, args ...interface{}) error {
	result, err := db.Exec(statement, args...)
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

type user struct {
	id        int
	firstName string
	lastName  string
}

func (u *user) String() string {
	return fmt.Sprintf("id: %d, first_name: %s, last_name: %s", u.id, u.firstName, u.lastName)
}

func selectUsersRecords(db database, query string) ([]*user, error) {
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

type status struct {
	id           int
	firstStatus  string
	secondStatus string
}

func (u *status) String() string {
	return fmt.Sprintf("id: %d, first_status: %s, second_status: %s", u.id, u.firstStatus, u.secondStatus)
}

func selectStatusRecords(db database, query string) ([]*status, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := make([]*status, 0)

	for rows.Next() {
		s := &status{}
		if err := rows.Scan(&s.id, &s.firstStatus, &s.secondStatus); err != nil {
			return nil, err
		}

		statuses = append(statuses, s)
	}

	if !rows.NextResultSet() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return statuses, nil
}
