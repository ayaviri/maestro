# README

This is a collection of services that compose the backend of `maestro`

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
# go run db/migrate.go [migration_name] [up/down]
# Example usage:
$ go run db/migrate.go create_user_table up
```

 **IMPORTANT: Ensure that both up/down scripts are idempotent. The go script that executes these does NOT rollback changes if left in an intermediate state**
