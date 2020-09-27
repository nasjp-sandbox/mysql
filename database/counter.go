package database

import "fmt"

type Counter struct {
	ID    int
	Count int
}

func (u *Counter) String() string {
	return fmt.Sprintf("id: %d, count: %d", u.ID, u.Count)
}

func SelectUsersRecords(db database, query string) ([]*Counter, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*Counter, 0)

	for rows.Next() {
		u := &Counter{}
		if err := rows.Scan(&u.ID, &u.Count); err != nil {
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
