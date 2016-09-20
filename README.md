# go-tx-guard

Simple database transaction guard

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


