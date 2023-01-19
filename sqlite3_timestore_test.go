package sqlite3

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

// TestTimestore tests that the _time_format config option works as expected.
// This would be moved into sqlite3_test.go if merged into the upstream repo.
func TestTimestore(t *testing.T) {
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

		// Create table. The time column will be sn integer type so it isn't automatically
		// converted to a time.Time on read. Need to introspect what was written.
		_, err = db.Exec("CREATE TABLE foo(ts INTEGER)")
		if err != nil {
			t.Fatal("Failed to create table:", err)
		}

		ts := time.Date(2023, 1, 19, 13, 45, 35, 45028023, time.UTC)
		_, err = db.Exec("INSERT INTO foo(ts) VALUES(?)", ts)
		if err != nil {
			t.Fatal("Failed to insert timestamp:", err)
		}

		// Test that time is stored correctly
		var tsInt int64
		var tsStr string
		row := db.QueryRow("SELECT ts FROM foo LIMIT 1;")

		if mode == "unix" || mode == "unix_ms" {
			err = row.Scan(&tsInt)
		} else {
			err = row.Scan(&tsStr)
		}

		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
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
		switch err := row.Scan(&tsTime); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			if tsTime != ts {
				t.Errorf("Timestamp value should be %v, not %v", ts, tsTime)
			}
		}
	}
}
