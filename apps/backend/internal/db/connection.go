package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"maestro/internal"
)

func EstablishConnection(dbPtr **sql.DB) {
	*dbPtr, err = sql.Open("sqlite3", "./db/maestro.db")
	internal.HandleError(err, "Failed to connect to database")
	err = (*dbPtr).Ping()
	internal.HandleError(err, "Could not ping database")
}
