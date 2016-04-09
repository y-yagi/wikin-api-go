package main

import (
	"database/sql"
	"time"
)

type Page struct {
	Id         int           `json:"id"`
	Title      string        `json:"title"`
	Body       string        `json:"body"`
	Parent_Id  sql.NullInt64 `json:"parent_id"`
	Created_at time.Time     `json:"created_at"`
	Updated_at time.Time     `json:"updated_at"`
}
