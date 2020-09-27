package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
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

	errCh := make(chan error)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		tx1, err := db.Begin()
		if err != nil {
			errCh <- err
		}

		defer func() {
			if err := tx1.Rollback(); err != nil {
				if errors.Is(err, sql.ErrTxDone) {
					return
				}

				if err != nil {
					errCh <- err
				}
			}
		}()

		if err := exec(tx1, "UPDATE app.users SET first_name='jiro', last_name='sato' WHERE id=1"); err != nil {
			errCh <- err
		}

		if err := tx1.Commit(); err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		tx2, err := db.Begin()
		if err != nil {
			errCh <- err
		}

		defer func() {
			err := tx2.Rollback()
			if errors.Is(err, sql.ErrTxDone) {
				return
			}

			if err != nil {
				errCh <- err
			}
		}()

		if err := exec(tx2, "UPDATE app.users SET first_name='saburo', last_name='kobayashi' WHERE id=1"); err != nil {
			errCh <- err
		}

		if err := tx2.Commit(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	users, err := selectUsersRecords(db, "SELECT * FROM app.users")
	if err != nil {
		return err
	}

	fmt.Println(users)
	// [id: 1, first_name: saburo, last_name: kobayashi]
	// [id: 1, first_name: jiro, last_name: sato]
	// どちらになるかわからない
	// 更新のロスト

	return nil
}
