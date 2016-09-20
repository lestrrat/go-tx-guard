# go-tx-guard

Simple database transaction guard

[![Build Status](https://travis-ci.org/lestrrat/go-tx-guard.png?branch=master)](https://travis-ci.org/lestrrat/go-tx-guard)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-tx-guard?status.svg)](https://godoc.org/github.com/lestrrat/go-tx-guard)

# SYNOPSIS

```go
import (
  "database/sql"

  "github.com/lestrrat/go-tx-guard"
)

func main() {
  db, err := guard.Open("mysql", "....")
  if err != nil {
    println(err.Error())
    return
  }

  tx, err := db.Begin()
  if err != nil {
    println(err.Error())
    return
  }
  // if tx.Commit or tx.Rollback is never explicitly
  // called, the transaction is automatically rolled back
  defer tx.AutoRollback()

  // do stuff with tx, maybe insert and then commit
}
```

# DESCRIPTION

Often times we would like to have an automatic rollback in case we return
early from a method performing an SQL transaction operation. This small
wrapper creates transaction objects that have an `AutoRollback` method that
you can safely register to a `defer` statement so that in case of errors
`Rollback` gets called.
