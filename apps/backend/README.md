# README

This is a collection of services that compose the backend of `maestro`

## Architectural Diagram

Below are the primary components of the backend, and the arrows describe the "checkout" process, in which a user downloads all of the songs they've added to their "cart"

```
                                  ┌──────────────────────────┐                                          
                                  │                          │                                          
                                  │                          │                                          
                                  │   frontend application   │                                          
                                  │                          │                                          
                                  │                          │                                          
                                  └───┬───────┬──────┬───────┘                                          
                                      │ ▲     │ ▲    │ ▲                                                
                                    1.│ │5.   │ │    │ │                                                
                                      │ │     │ │    │ │                                                
                                      │ │   6.│ │11. │ │                                                
                                      │ │     │ │    │ │                                                
                                      │ │     │ │    │ │                                                
frontend                              │ │     │ │    │ │                                                
──────────────────────────────────────┼─┼─────┼─┼────┼─┼─────────────────────────────────────────────   
backend                               │ │     │ │ 12.│ │14.                                             
                                      ▼ │     ▼ │    ▼ │                                                
                                    ┌───┴───────┴──────┴───┐                                            
                                    │                      │                                            
                                    │                      │                                            
                 1.5.               │                      │                                            
             ┌──────────────────────┼   core web server    │                                            
             │ ┌────────────────────►                      │                                            
             │ │     2.             │                      │                                            
             │ │                    │                      │                                            
             │ │                    └──┬───────────────────┘                                            
             │ │                       │        ▲       ▲                 worker environment            
             ▼ │                       │        │       │                 ┌─────────────────────────┐   
  ┌────────────┴───────────────┐       │        │       │                 │                         │   
  │                            │       │        │       │                 │                         │   
  │                            │       │        │       │    13.          │    ┌───────────────┐    │   
  │ postgresql database server │       │        │       └─────────────────┼────►               │    │   
  │                            │     3.│   7/10.│                         │    │  file server  │    │   
  │                            │       │        │                         │    │               │    │   
  └────────────────────────────┘       │        │                         │    └───────────────┘    │   
                                       │        │                         │                         │   
                                       │        │                         │                         │   
                                       │        │                         │      ┌────────────┐     │   
 ┌───────────────────────────┐         │        │                         │      │            │     │   
 │                           │◄────────┘        │      4.                 │      │   worker   │     │   
 │   message queue server    ├──────────────────┼─────────────────────────┼─────►│            │     │   
 │                           │                  │                         │      └────┬─▲─────┘     │   
 └────────────────────────▲──┘                  │                         │           │ │           │   
                      ▲   │                     │                         │           │ │ 8.        │   
                      │   └─────────────────────┘                         │           │ │           │   
                      │                                                   └───────────┼─┼───────────┘   
                      │          9.                                                   │ │               
                      └───────────────────────────────────────────────────────────────┘ │               
                                                                                        │               
 ───────────────────────────────────────────────────────────────────────────────────────┼───────────────
                                                                                        │               
                                                                                        │               
                                        ┌─────────────────┐                             │               
                                        │                 │◄────────────────────────────┘               
                                        │     youtube     │                                             
                                        │                 │                                             
                                        └─────────────────┘                                             
```
1. Frontend application sends `/checkout` request
1.5. Core web server authenticates user by comparing request bearer token to token in Postgres
2. Core web server obtains cart items from Postgres
3. Core web server posts message to the message queue (RabbitMQ server). This states the intent for a cart's worth of items to be downloaded
4. Worker picks up message asynchronously, begins download
5. ID of the download job is returned to frontend from the core web server
6. Frontend opens Server Sent Events (SSE) connection with core web server at `/job/{path}`
7. Core web server awaits message from queue signaling *completion* of checkout
8. Worker completes download
9. Worker writes message to queue signaling completion of checkout. Message contains download URL of songs, served by the file server in the same environment
10. Core web server picks up message of checkout completion
11. Core web server responds to frontend with download URLs
12. Frontend request download URLs, served by core web server at `/download/{path}`
13. Songs are downloaded from core web server at `/download/{path}`; core web server is serving as a proxy for the file server, in which the files live
14. Core web server responds with songs


## Project Structure

I'm new to Go, so documenting this for myself and others serves as a learning tool for me.

| Directory | Purpose |
| --------- | ------- |
| `cmd`     | Contains packages with entry points for invokation (eg. `cmd/core/main.go` contains the web server, `cmd/migrate/main.go` contains the DB migration script, `cmd/fs/main.go` contains the file server meant to be run as a separate process in the worker's environment) |
| `db`      | Contains the SQL migration scripts |
| `internal` | Contains utilities used throughout the project. Importable as `maestro/internal` |
| `worker`  | Contains the script the rabbitmq worker will run in order to consume checkout messages and download music out of band |

## Local Development

There is a docker-compose.yaml file to prop up the PostgreSQL and RabbitMQ servers, and there is a Bash script to prop up the core web server, the file server, and the worker process.
The docker-compose.yaml file expects certain environment variables, so there is a load_env.sh script for that purpose
```
$ source load_env.sh .env.core
$ docker-compose.yaml up -d
$ # Each in a separate terminal
$ ./dev.sh core
$ ./dev.sh worker
$ ./dev.sh fs
```

### Database Migrations

ASIDE: I tried to recreate an [Alembic](https://github.com/sqlalchemy/alembic) style migration system that emulates version control, but I realised that some operations could not be performed after table creation, so I decided to just throw everything in a single SQL script, essentially removing any value the Go migration script below has

In order to create migration:

1) Create a directory in the `db/migrate` directory (eg. `db/migrate/drop_user_table`). Name it appropriately after the migration being performed. 
2) In this directory, create two SQL scripts, `up.sql` and `down.sql`.
3) In `up.sql`, write the migration script
4) In `down.sql`, write a migration script that can undo the changes made in `up.sql`

To execute a migration, run the following

```
# Usage:
# go run cmd/migrate/main.go [migration_name] [up/down]
# Example usage:
$ go run cmd/migrate/main.go create_user_table up
```

 **IMPORTANT: Ensure that both up/down scripts are idempotent. The go script that executes these does NOT rollback changes if left in an intermediate state**
