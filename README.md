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

Normally `time.Time` values are stored as a string in the format: `2006-01-02 15:04:05.999999999-07:00`.
The `_time_format` connection option tells the driver to store times in an integral format:

- `unix` stores the time as a [Unix timestamp](https://en.wikipedia.org/wiki/Unix_time)
- `unix_ms` stores the time as a millisecond-precision Unix timestamp (i.e. integer number of milliseconds since the epoch)

Note that "unix milliseconds" is not a SQlite-supported type so you wont be able to use its built in time functions on the data directly without converting it. The data is parsable by the driver, however.

## Why?

I've always preferred the simplicity, compactness, and unambiguity of Unix timestamps. I also like using `time.Time` without having to add special handling or distinct types to effect integer storage. Modifying the driver wasn't my first choice but I think it's the simplest.

## Misc

This isn't related to the driver but is something I sometimes have to Google. If you want to default a column to the current unix timestamp you can use one of the following:

```sql
CREATE TABLE foo(ts INTEGER DEFAULT (strftime('%s', 'now')));
```

or if using SQLite 3.38.0 (2/2/2022) or later you can also do:

```sql
CREATE TABLE f3(ts INTEGER DEFAULT (unixepoch()));
```
