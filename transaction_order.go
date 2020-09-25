package main

import (
	"fmt"
)

func transactionOrder() error {
	db, err := connect()
	if err != nil {
		return err
	}

	if err := exec(db, "CREATE TABLE IF NOT EXISTS app.users (id int, first_name varchar(10), last_name varchar(10))"); err != nil {
		return err
	}

	defer exec(db, "DROP TABLE IF EXISTS app.users")

	if err := exec(db, "INSERT INTO app.users (id, first_name, last_name) VALUES (1, 'taro', 'tanaka')"); err != nil {
		return err
	}

	tx1, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx1.Rollback()

	tx2, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx2.Rollback()

	errCh := make(chan error)

	// こっちのほうが後
	go func() {
		defer close(errCh)

		if err := exec(tx2, "UPDATE app.users SET first_name='saburo', last_name='kobayashi' WHERE id=1"); err != nil {
			errCh <- err
		}

		if err := tx2.Commit(); err != nil {
			errCh <- err
		}
	}()

	// こっちが必ず先
	// transactionが貼られたのが先だから
	if err := exec(tx1, "UPDATE app.users SET first_name='jiro', last_name='sato' WHERE id=1"); err != nil {
		return err
	}

	if err := tx1.Commit(); err != nil {
		return err
	}

	if err := <-errCh; err != nil {
		return err
	}

	users, err := selectUsersRecords(db, "SELECT * FROM app.users")
	if err != nil {
		return err
	}

	fmt.Println(users)
	// [id: 1, first_name: saburo, last_name: kobayashi]

	return nil
}
