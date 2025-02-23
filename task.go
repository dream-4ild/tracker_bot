package main

import (
	"database/sql"
	"time"
)

type Task struct {
	id        int64
	userId    int64
	project   string
	task      string
	status    string
	deadline  sql.NullTime
	updatedAt time.Time
}
