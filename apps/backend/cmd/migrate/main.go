package main

import (
	"database/sql"
	"errors"
	"fmt"
	"maestro/internal"
	xdb "maestro/internal/db"
	"os"
	"path"

	"github.com/ayaviri/goutils/timer"
)

var err error

// Reads and validates the command line arguments. Reads the
// desired SQL migration script into memory and returns it, along
// with any encountered error
func readSqlScript() (string, error) {
	var args []string = os.Args[1:] // Drops the program name

	if len(args) != 2 {
		return "", errors.New("Must provider migration name and migration direction")
	}

	var migrationDirectory string = args[0]
	var migrationDirection string = args[1]

	if migrationDirection != "up" && migrationDirection != "down" {
		return "", errors.New(
			"Invalid migration direction received. Must be one of 'up' or 'down'",
		)
	}

	var scriptContents []byte
	scriptContents, err = os.ReadFile(
		path.Join(
			"./db/migrate",
			migrationDirectory,
			fmt.Sprintf("%s.sql", migrationDirection),
		),
	)

	if err != nil {
		return "", fmt.Errorf("Could not read contents of script: %v\n", err.Error())
	}

	return string(scriptContents), nil
}

// Usage: go run db/migrate.go [name_of_migration] [up/down]
// Example: go run db/migrate.go create_initial_tables up
func main() {
	var script string

	timer.WithTimer("pulling script from disk", func() {
		script, err = readSqlScript()
		internal.HandleError(
			err,
			"Failed to parse command line arguments and read migration script",
		)
	})

	var db *sql.DB

	timer.WithTimer("connecting to database", func() {
		xdb.EstablishConnection(&db)
	})

	defer db.Close()

	timer.WithTimer("executing migration script", func() {
		_, err = db.Exec(script)
		// TODO: This might leave the DB in an intermediate state, so it's
		// important to make both migration directions idempotent
		internal.HandleError(err, "Failed to execute migration")
	})
}
