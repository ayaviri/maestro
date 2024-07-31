package main

import (
	"database/sql"
	"errors"
	"fmt"
	"maestro/utils"
	"os"
	"path"
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

	return string(scriptContents), nil
}

// Usage: go run db/migrate.go [name_of_migration] [up/down]
// Example: go run db/migrate.go create_initial_tables up
func main() {
	var script string

	utils.WithTimer("pulling script from disk", func() {
		script, err = readSqlScript()
		utils.HandleError(
			err,
			"Failed to parse command line arguments and read migration script",
		)
	})

	var db *sql.DB

	utils.WithTimer("connecting to database", func() {
		db, err = sql.Open("", "./db/maestro.db")
		utils.HandleError(err, "Failed to connect to database")
	})

	defer db.Close()

}
