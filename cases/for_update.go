package cases

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nasjp-sandbox/mysql/database"
	"golang.org/x/sync/errgroup"
)

func ForUpdate() error {
	fmt.Println("start: for update")

	db, err := database.Connect()
	if err != nil {
		return err
	}

	if err := database.Exec(db, "CREATE TABLE IF NOT EXISTS app.statuses (id int, first_status tinyint(1), second_status tinyint(1))"); err != nil {
		return err
	}

	defer database.Exec(db, "DROP TABLE IF EXISTS app.statuses")

	if err := database.Exec(db, "INSERT INTO app.statuses (id, first_status, second_status) VALUES (1, 0, 0)"); err != nil {
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

	g := new(errgroup.Group)

	g.Go(func() error {
		var rerr error

		defer func() {
			err := tx1.Rollback()
			if errors.Is(err, sql.ErrTxDone) {
				return
			}

			if err != nil {
				rerr = err
			}
		}()

		// statuses, err := database.SelectStatusRecords(tx1, "SELECT * FROM app.statuses WHERE id=1")
		statuses, err := database.SelectStatusRecords(tx1, "SELECT * FROM app.statuses WHERE id=1 FOR UPDATE")
		if err != nil {
			return err
		}

		if err := database.Exec(tx1, "UPDATE app.statuses SET first_status=1, second_status=? WHERE id=1", statuses[0].SecondStatus); err != nil {
			return err
		}

		if err := tx1.Commit(); err != nil {
			return err
		}

		return rerr
	})

	g.Go(func() error {
		var rerr error

		defer func() {
			err := tx2.Rollback()
			if errors.Is(err, sql.ErrTxDone) {
				return
			}

			if err != nil {
				rerr = err
			}
		}()

		// statuses, err := database.SelectStatusRecords(tx2, "SELECT * FROM app.statuses WHERE id=1")
		statuses, err := database.SelectStatusRecords(tx2, "SELECT * FROM app.statuses WHERE id=1 FOR UPDATE")
		if err != nil {
			return err
		}

		if err := database.Exec(tx2, "UPDATE app.statuses SET first_status=?, second_status=1 WHERE id=1", statuses[0].FirstStatus); err != nil {
			return err
		}

		if err := tx2.Commit(); err != nil {
			return err
		}

		return rerr
	})

	if err := g.Wait(); err != nil {
		return err
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
