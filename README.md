# go-sqlite3

This is a fork of [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) with the following changes:

- Support for a new `_time_format` connection option

I will update the fork from the upstream repo as needed, or when requested. If this option turns out to be
useful I'll consider submitting a PR upstream.

## Usage

This is a new option for the DSN:

```go
db, err := sql.Open("sqlite3", "filename.db?_time_format=<format>")
```

The `_time_format` connection option tells the driver to store times in an integral format instead of the lengthy (30+ char) ISO8601-style strings:

- `unix` stores the time as a [Unix timestamp](https://en.wikipedia.org/wiki/Unix_time)
- `unix_ms` stores the time as a millisecond-precision Unix timestamp (i.e. integer number of milliseconds since the epoch)

Note:

- Even though the underlying storage is integer, the column time needs to be one of: `date`, `datetime`, `timestamp`. These are the column types `mattn/go-sqlite` will convert into `time.Time` (no change from the original driver).
- "unix milliseconds" are not a SQlite-supported type so you won't be able to use its built-in time functions on the data directly without converting it. The format has been parsable by the driver for many years, however.

## Why?

I've always preferred the simplicity, compactness, and unambiguity of Unix timestamps. I also like using `time.Time` without having to add special handling or distinct types to effect integer storage. Modifying the driver wasn't my first choice but I think it's the simplest.
