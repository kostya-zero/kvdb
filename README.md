# `KVDB`

KVDB is a lightweight key‑value store written in Go. It keeps data in
named maps so multiple services can share a single server without clobbering
each other's keys. The server exposes a plain TCP protocol and can optionally
persist all updates to a file.

## Installation

#### GitHub Releases

Download an archive from [GitHub Releases]("https://github.com/kostya-zero/kvdb/releases") and extract the KVDB binary
to directory that is added to `PATH`.

#### Docker

Clone this repository and use Docker CLI to build and run container.

```shell
docker build -t kvdb .
docker run -p 5511:5511 kvdb
```

## Usage

Run the server with the `serve` command.

```shell
kvdb serve
```

All options for KVDB should be provided as environment variables. Available variables:

- `KVDB_PORT` - TCP port to bind to (default `5511`).
- `KVDB_DATABASE` - path to the database file. When omitted, data is kept in memory only.
- `KVDB_SAVE_INTERVAL` - the save interval in milliseconds. After this time, KVDB will run save if it is needed.

Examples:

```shell
# Start with defaults (port 5511, in-memory database)
kvdb serve

# Persist data to kvdb.db and expose on custom port
export KVDB_PORT=7777
export KVDB_DATABASE/var/lib/kvdb/kvdb.db
kvdb serve
```

## API overview

The server communicates over plain TCP.
Each request is a single line command and the response is a text string.
The basic commands are:

- `CREATEDB <db>` – create a new database map.
- `REMOVE DB <db>` – remove a database.
- `REMOVE KEY <db>.<key>` – delete a key.
- `SET <db>.<key> "<value>"` – add a new key with value.
- `GET <db>.<key>` – return the value for a key.
- `UPDATE <db>.<key> "<value>"` – replace the current value.
- `LIST <db>` - print list of databases (if database name is not provided) or keys inside database.

Responses are either `OK` or one of the following error codes: `ALREADY_EXISTS`,
`DATABASE_IS_EMPTY`, `DATABASE_NOT_FOUND`, `KEY_NOT_FOUND`,`KEY_NOT_PROVIDED`,
`ILLEGAL_CHARACTERS`, or `BAD_QUERY`.

You can test the server with tools like `nc`:

```shell
$ nc 127.0.0.1 5511
CREATEDB users
SET users.alice "alice@example.com"
GET users.alice
```

## License

KVDB is licensed under MIT License. Learn more in [LICENSE](LICENSE) file.
