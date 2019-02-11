package sqlite3

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func CreateSqlite3Pool(m map[string]interface{}) (BasePool, error) {
	return CreatePool(&PoolConfig{
		Min: m["min"].(int),
		Max: m["max"].(int),
		Create: func() (interface{}, error) {
			db, err := sql.Open("sqlite3", m["db"].(string))
			if err != nil {
				return nil, err
			}
			conn := &Sqlite3Connect{
				db,
			}
			return conn, nil
		},
		Close: func(v interface{}) error {
			return v.(Sqlite3Connect).Db.Close()
		},
		IdleTimeout: 15 * time.Second,
	})
}
