package database

import "fmt"

type User struct {
	ID        int
	FirstName string
	LastName  string
}

func (u *User) String() string {
	return fmt.Sprintf("id: %d, first_name: %s, last_name: %s", u.ID, u.FirstName, u.LastName)
}

func SelectUsersRecords(db database, query string) ([]*User, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)

	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName); err != nil {
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
