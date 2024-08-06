package db

import (
	"database/sql"
)

// Expects a raw SQL count query, writing the first column of the first row
// returned into a string, converting it into an integer, and returning it.
// Returns any errors encountered along the way as well
func QueryCount(db *sql.DB, query string) (int64, error) {
	var count int64
	var row *sql.Row = db.QueryRow(query)
	err = row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
