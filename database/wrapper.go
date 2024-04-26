package database

import "database/sql"

type Querier interface {
	Exec(query string, args ...any) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
}
