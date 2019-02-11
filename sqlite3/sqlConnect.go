package sqlite3

import "database/sql"

type SqlConnect interface {
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Query(sql string, args ...interface{}) (*sql.Rows, error)
}
type Sqlite3Connect struct {
	Db *sql.DB
}

func (connect Sqlite3Connect) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return connect.Db.Exec(sql, args...)
}

func (connect Sqlite3Connect) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return connect.Db.Query(sql, args...)
}
