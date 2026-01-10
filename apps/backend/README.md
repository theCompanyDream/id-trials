# Backend Service

A Go-based backend service built with **Cobra** for CLI ergonomics and **Echo** for HTTP routing. It compiles into a single binary that can either **generate data** or **run an HTTP server**, depending on the command you invoke.

Think of it as a two-mode creature: one part data factory, one part web server.

## Tech Stack

* **Go**
* **Cobra** – CLI command and flag management
* **Echo** – HTTP server framework
* **dotenv-compatible configuration** via environment variables

## Configuration

The service reads configuration from environment variables. You can define them directly in your shell or via a `.env` file.

Common variables include:

```env
BACKEND_HOST=localhost
BACKEND_PORT=8080
```

Additional variables required for record generation should also live in `.env`.


## Commands

### 1. Generate Records

Generates a batch of records using values from the `.env` file.

```bash
./backend generate --records 10000 --batch 100
```

#### Flags

| Flag        | Description                         | Required |
| -- | -- | -- |
| `--records` | Total number of records to generate | Yes      |
| `--batch`   | Number of records per batch         | Yes      |

This command is useful for seeding databases, load testing, or feeding hungry downstream systems.



### 2. Run Server

Starts the HTTP server.

```bash
./backend server --host localhost --port 234
```

#### Flags

| Flag     | Description       | Environment Fallback |
| -- | -- | -- |
| `--host` | Host to bind to   | `BACKEND_HOST`       |
| `--port` | Port to listen on | `BACKEND_PORT`       |

If a flag is omitted, the corresponding environment variable is used automatically.

Example with environment variables only:

```bash
./backend server
```



## Environment Variable Precedence

For the `server` command:

1. CLI flags take priority
2. Environment variables are used if flags are not provided
3. If neither is set, startup will fail (unless defaults are defined in code)



## Development

Run directly without building:

```bash
go run main.go generate --records 1000 --batch 50
go run main.go server
```

## Summary

* One binary
* Two commands
* Zero ambiguity about where configuration comes from

If something doesn’t start, the first place to look is your `.env`. The second place is the flags. After that, it’s probably intentional.

If you want, I can also:

* Add a **Makefile**
* Include **example `.env`**
* Tighten this up into a more “production-grade” README with logging, health checks, and exit codes
