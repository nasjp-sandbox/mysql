package cases

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/nasjp-sandbox/mysql/database"
)

func AtomicWrite() error {
	fmt.Println("start: atomic write")

	db, err := database.Connect()
	if err != nil {
		return err
	}

	if err := database.Exec(db, "CREATE TABLE IF NOT EXISTS app.counters (id int, count int)"); err != nil {
		return err
	}

	defer database.Exec(db, "DROP TABLE IF EXISTS app.counters")

	if err := database.Exec(db, "INSERT INTO app.counters (id, count) VALUES (1, 0)"); err != nil {
		return err
	}

	tx1, err := db.Begin()
	if err != nil {
		return err
	}

	tx2, err := db.Begin()
	if err != nil {
		return err
	}

	errCh := make(chan error)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

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

		if err := database.Exec(tx1, "UPDATE app.counters SET count=count+1 WHERE id=1"); err != nil {
			errCh <- err
		}

		if err := tx1.Commit(); err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		defer func() {
			err := tx2.Rollback()
			if errors.Is(err, sql.ErrTxDone) {
				return
			}

			if err != nil {
				errCh <- err
			}
		}()

		if err := database.Exec(tx2, "UPDATE app.counters SET count=count+1 WHERE id=1"); err != nil {
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

	counters, err := database.SelectUsersRecords(db, "SELECT * FROM app.counters")
	if err != nil {
		return err
	}

	fmt.Println(counters)
	// [id: 1, count: 2]
	// となる

	// [id: 1, count: 1] のように更新のロストは起こらない
	// アトミックな更新処理

	return nil
}
