package db

import (
	"database/sql"
	"fmt"
	"maestro/internal"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func getDatabaseUrl() string {
	err = godotenv.Load()
	internal.HandleError(err, "Could not load environment variables")
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
	*dbPtr, err = sql.Open("pgx", getDatabaseUrl())
	internal.HandleError(err, "Failed to connect to database")
	err = (*dbPtr).Ping()
	internal.HandleError(err, "Could not ping database")
}
