package cases

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/nasjp-sandbox/mysql/database"
)

func ForUpdate() error {
	fmt.Println("start: for update")

	db, err := database.Connect()
	if err != nil {
		return err
	}

	if err := database.Exec(db, "CREATE TABLE IF NOT EXISTS app.statuses (id int, first_status tinyint(1) NOT NULL, second_status tinyint(1) NOT NULL)"); err != nil {
		return err
	}

	defer database.Exec(db, "DROP TABLE IF EXISTS app.statuses")

	if err := database.Exec(db, "INSERT INTO app.statuses (id, first_status, second_status) VALUES (1, 0, 0)"); err != nil {
		return err
	}

	errCh := make(chan error)

	var wg sync.WaitGroup

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
		// statuses, err := database.SelectStatusRecords(tx2, "SELECT * FROM app.statuses WHERE id=1")
		statuses, err := database.SelectStatusRecords(tx2, "SELECT * FROM app.statuses WHERE id=1 FOR UPDATE")
		if err != nil {
			errCh <- err
		}

		if err := database.Exec(tx2, "UPDATE app.statuses SET first_status=?, second_status=1 WHERE id=1", statuses[0].FirstStatus); err != nil {
			errCh <- err
		}

		if err := tx2.Commit(); err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		tx1, err := db.Begin()
		if err != nil {
			errCh <- err
		}

		defer func() {
			err := tx1.Rollback()
			if errors.Is(err, sql.ErrTxDone) {
				return
			}

			if err != nil {
				errCh <- err
			}
		}()

		// statuses, err := database.SelectStatusRecords(tx1, "SELECT * FROM app.statuses WHERE id=1")
		statuses, err := database.SelectStatusRecords(tx1, "SELECT * FROM app.statuses WHERE id=1 FOR UPDATE")
		if err != nil {
			errCh <- err
		}

		if err := database.Exec(tx1, "UPDATE app.statuses SET first_status=1, second_status=? WHERE id=1", statuses[0].SecondStatus); err != nil {
			errCh <- err
		}

		if err := tx1.Commit(); err != nil {
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

	results, err := database.SelectStatusRecords(db, "SELECT * FROM app.statuses")
	if err != nil {
		return err
	}

	fmt.Println(results)
	// [id: 1, first_status: 1, last_status: 1]
	// for updateで行ロックをしないと
	// [id: 1, first_status: 1, last_status: 0] もしくは
	// [id: 1, first_status: 0, last_status: 1] になってしまう

	return nil
}
