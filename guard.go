package guard

import (
	"database/sql"
	"fmt"
	"sync"

	gsg "github.com/lestrrat/go-simple-guard"
)

/*
package guard implements a simple transaction guard around the
default database/sql library. You just need to wrap the default
database/sql handle:

	func Connect() *guard.DB {
		realDB, err := sql.Open(driverName, dsn)
		if err != nil {
			...
		}
		return &guard.DB{realDB}
	}

Then after this, when you use transactions you can use the following
construct:

	func Foo(db *guard.DB) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		defer tx.AutoRollback() // Automatically rollback upon failure to commit

		.... // Code that might fail
		if err != nil {
			// AutoRollback kicks in, and ROLLBACK is issued.
			return err
		}

		// If this commit is successful, AutoRollback doesn't do anything
		tx.Commit()
	}

*/

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
	lock        sync.Locker
	Name        string
	afterCommit []func()
}

// Begin begins a transactin, and creates a new Tx object.
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	ttx := &Tx{
		lock: &sync.Mutex{},
		gcb: gsg.NewCB(func() error {
			return tx.Rollback()
		}),
		Tx: tx,
	}
	ttx.Name = fmt.Sprintf("%p", ttx)
	return ttx, nil
}

// Commit sets the finished flag and then calls Commit() on the underlying
// sql.Tx object. Failure to Commit() by errors do not affect the finished
// flag being set. After calling this method, AutoRollback() is a no op
func (tx *Tx) Commit() error {
	defer tx.gcb.Cancel()
	if err := tx.Tx.Commit(); err != nil {
		return err
	}
	defer tx.runAfterCommit()
	return nil
}

// AddAfterCommit adds a callback that gets called only when
// a Commit() was successful. The callbacks are executed in
// the order that they were added. Errors are ignored.
// If a panic occurs within one of these callbacks, execution
// of the callbacks stop there.
func (tx *Tx) AddAfterCommit(cb func()) {
	tx.afterCommit = append(tx.afterCommit, cb)
}

// AfterCommit hooks do NOT report errors! be careful
func (tx *Tx) runAfterCommit() {
	defer recover()
	for _, cb := range tx.afterCommit {
		cb()
	}
}

// Rollback sets the finished flag and then calls Rollback() on the underlying
// sql.Tx object. Failure to Rollback() by errors do not affect the finished
// flag being set. After calling this method, AutoRollback() is a no op
func (tx *Tx) Rollback() error {
	defer tx.gcb.Cancel()
	return tx.Tx.Rollback()
}

// AutoRollback only rollsback if Commit/Rollback has not been
// previously called. This way you can safely call
//
//     tx, err := db.Begin()
//     defer tx.AutoRollback()
//
func (tx *Tx) AutoRollback() error {
	return tx.gcb.Fire()
}