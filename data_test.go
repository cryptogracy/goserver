package main

import (
  "testing"
)

func TestDBInit(t *testing.T) {
  db, err := DBInit(":memory:")
  if err != nil { t.Error(err)}

	var res string
	err = db.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='Files'").Scan(&res)
  if err != nil { t.Error(err) }
}
