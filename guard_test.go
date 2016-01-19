package guard

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// This is a silly test that doesn't do anything...
func TestSQLite3(t *testing.T) {
	rdb, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("Failed to open database: %s", err)
		return
	}
	db := &DB{rdb}

	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to start transaction: %s", err)
		return
	}
	defer tx.AutoRollback()

	afterCommitCalled := 0
	tx.AddAfterCommit(func() {
		afterCommitCalled++
	})

	tx.Commit()

	if afterCommitCalled != 1 {
		t.Errorf("Expected AfterCommit hook to be called once, got %d", afterCommitCalled)
	}
}
