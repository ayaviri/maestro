package db

import (
	"database/sql"
	"fmt"
	"maestro/internal"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func getDatabaseServerUrl() string {
	databaseUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("POSTGRES_DB"),
	)

	return databaseUrl
}

func EstablishConnection(dbPtr **sql.DB) {
	var url string = getDatabaseServerUrl()
	*dbPtr, err = sql.Open("pgx", url)
	internal.HandleError(
		err,
		fmt.Sprintf("Failed to connect to database, given URL: %s\n", url),
	)
	err = (*dbPtr).Ping()
	internal.HandleError(err, "Could not ping database")
}
