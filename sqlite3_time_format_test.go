package sqlite3

import (
	"database/sql"
	"testing"
	"time"
)

// TestTimeFormat tests that the _time_format config option works as expected.
// This would be moved into sqlite3_test.go if merged into the upstream repo.
func TestTimeFormat(t *testing.T) {
	modes := []string{"", "unix", "unix_ms", "nope"}
	for _, mode := range modes {
		db, err := sql.Open("sqlite3", ":memory:?_time_format="+mode)
		if err != nil {
			t.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		// A query is required to test for errors on the config
		_, err = db.Exec("SELECT 1")
		switch mode {
		// Valid modes
		case "unix", "unix_ms", "":
			if err != nil {
				t.Fatal("Failed to open database:", err)
			}
		default:
			if err == nil {
				t.Fatal("Expected error for invalid time format")
			}
			continue
		}

		// For the test both a time-style type (timestamp, date or datetime) and and integer
		// are used. The first type is the normal use case and the driver will automatically
		// returns a time.Time. But an integer is also used for the test so we can see what
		// was really written. (Note: you can't use an integer column and expect to scan back
		// into time.Time... it will error out)
		_, err = db.Exec("CREATE TABLE foo(ts TIMESTAMP, ts_int INTEGER)")
		if err != nil {
			t.Fatal("Failed to create table:", err)
		}

		ts := time.Date(2023, 1, 19, 13, 45, 35, 45028023, time.UTC)
		_, err = db.Exec("INSERT INTO foo(ts, ts_int) VALUES(?, ?)", ts, ts)
		if err != nil {
			t.Fatal("Failed to insert timestamp:", err)
		}

		// Test that time is stored correctly
		var tsInt int64
		var tsStr string
		row := db.QueryRow("SELECT ts_int FROM foo LIMIT 1;")

		if mode == "unix" || mode == "unix_ms" {
			err = row.Scan(&tsInt)
		} else {
			err = row.Scan(&tsStr)
		}

		switch err {
		case nil:
			if mode == "unix" {
				if tsInt != ts.Unix() {
					t.Errorf("Timestamp value should be %v, not %v", ts.Unix(), tsInt)
				}
			} else if mode == "unix_ms" {
				if tsInt != ts.UnixMilli() {
					t.Errorf("Timestamp value should be %v, not %v", ts.UnixMilli(), tsInt)
				}
			} else {
				if tsStr != ts.Format(SQLiteTimestampFormats[0]) {
					t.Errorf("Timestamp value should be %v, not %v", ts.Format(SQLiteTimestampFormats[0]), tsStr)
				}
			}
		default:
			t.Fatal("error reading row:", err)
		}

		// Test that reading into time.Time works
		var tsTime time.Time
		row = db.QueryRow("SELECT ts FROM foo LIMIT 1;")
		if err := row.Scan(&tsTime); err != nil {
			t.Fatal("unexpected error:", err)
		}
		expTs := ts
		switch mode {
		case "unix":
			expTs = ts.Truncate(time.Second)
		case "unix_ms":
			expTs = ts.Truncate(time.Millisecond)
		}

		if tsTime != expTs {
			t.Errorf("Timestamp value should be %v, not %v", ts, tsTime)
		}
	}
}
