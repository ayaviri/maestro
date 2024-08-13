# README

This is a collection of services that compose the backend of `maestro`

## Project Structure

I'm new to Go, so documenting this for myself and others serves as a learning tool for me.

| Directory | Purpose |
| --------- | ------- |
| `cmd`     | Contains packages with entry points for invokation (eg. `cmd/core/main.go` contains the web server, `cmd/migrate/main.go` contains the DB migration script, `cmd/fs/main.go` contains the file server meant to be run as a separate process in the worker's environment) |
| `db`      | Contains the SQLite database and the SQL migration scripts |
| `internal` | Contains utilities used throughout the project. Importable as `maestro/internal` |
| `worker`  | Contains the script the rabbitmq worker will run in order to consume checkout messages and download music out of band |

## Database

`sqlite3` is used. For local development, run the following to prop up the database

```
$ touch db/maestro.db
```

### Migrations

1) Create a directory in the `db/migrate` directory (eg. `db/migrate/drop_user_table`). Name it appropriately after the migration being performed. 
2) In this directory, create two SQL scripts, `up.sql` and `down.sql`.
3) In `up.sql`, write the migration script
4) In `down.sql`, write a migration script that can undo the changes made in `up.sql`

To execute a migration, run the following

```
# Usage:
# go run cmd/migrate/main.go [migration_name] [up/down]
# Example usage:
$ go run db/migrate.go create_user_table up
```

 **IMPORTANT: Ensure that both up/down scripts are idempotent. The go script that executes these does NOT rollback changes if left in an intermediate state**
