package redis

import (
	"time"

	"github.com/go-redis/redis"
	Db "github.com/qiongtubao/latte_go_db"
)

func CreateRedisPool(m map[string]interface{}) (Db.BasePool, error) {
	return Db.CreatePool(&Db.PoolConfig{
		Min: m["min"].(int),
		Max: m["max"].(int),
		Create: func() (interface{}, error) {
			client := redis.NewClient(&redis.Options{
				Addr: m["host"].(string) + ":" + m["port"].(string),
				// Password: m["password"],
				DB: m["db"].(int),
			})
			conn := &RedisConnect{
				client,
			}
			return conn, nil
		},
		Close: func(v interface{}) error {
			v.(*RedisConnect).Close()
			return nil
		},
		IdleTimeout: 15 * time.Second,
	})
}
