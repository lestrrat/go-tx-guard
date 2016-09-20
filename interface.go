package guard

import (
	"database/sql"
	"sync"

	gsg "github.com/lestrrat/go-simple-guard"
)

// DB wraps a sql.DB object.
type DB struct {
	*sql.DB
}

// Tx wraps a sql.Tx object. The only difference between sql.Tx is that
// this has an AutoRollback method, which only calls Rollback() if the
// transaction hasn't already been committed or rolled back.
type Tx struct {
	*sql.Tx
	gcb         gsg.Guard
	mutex       sync.RWMutex
	Name        string
	afterCommit []func()
}
