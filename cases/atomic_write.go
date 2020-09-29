package cases

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nasjp-sandbox/mysql/database"
	"golang.org/x/sync/errgroup"
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

	g := new(errgroup.Group)

	g.Go(func() error {
		var rerr error

		defer func() {
			if err := tx1.Rollback(); err != nil {
				if errors.Is(err, sql.ErrTxDone) {
					return
				}

				if err != nil {
					rerr = err
				}
			}
		}()

		if err := database.Exec(tx1, "UPDATE app.counters SET count=count+1 WHERE id=1"); err != nil {
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

		if err := database.Exec(tx2, "UPDATE app.counters SET count=count+1 WHERE id=1"); err != nil {
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
