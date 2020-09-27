package database

import "fmt"

type Status struct {
	ID           int
	FirstStatus  string
	SecondStatus string
}

func (u *Status) String() string {
	return fmt.Sprintf("id: %d, first_status: %s, second_status: %s", u.ID, u.FirstStatus, u.SecondStatus)
}

func SelectStatusRecords(db database, query string) ([]*Status, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := make([]*Status, 0)

	for rows.Next() {
		s := &Status{}
		if err := rows.Scan(&s.ID, &s.FirstStatus, &s.SecondStatus); err != nil {
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
