package string

import (
	"fmt"
	"testing"

	Redis "github.com/qiongtubao/latte_go_db/redis"
)

func Test_Add(t *testing.T) {
	pool, err := Redis.CreateRedisPool(map[string]interface{}{
		"min":  1,
		"max":  2,
		"host": "localhost",
		"port": "6379",
		"db":   0,
	})
	if err != nil {
		t.Error(err)
	}
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
	}

	model := CreateModel("hello")
	model.Set("world")
	err = Update(conn.(*Redis.RedisConnect), model)
	if err != nil {
		t.Error(err)
	}
	// model := RedisString.Get(conn, "hello")
	// if model.Interface() == nil  {
	// 	RedisString.Add(conn, model)
	// }else {
	// 	fmt.Println(model.String())
	// }
	// model.Set("world")
	// RedisString.Update(conn, model)
}

func Test_Query(t *testing.T) {
	pool, err := Redis.CreateRedisPool(map[string]interface{}{
		"min":  1,
		"max":  2,
		"host": "localhost",
		"port": "6379",
		"db":   0,
	})
	if err != nil {
		t.Error(err)
	}
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
	}

	model, err := Get(conn.(*Redis.RedisConnect), "hello")
	if err != nil {
		t.Error(err)
	}
	if model.Interface() == nil {
		Update(conn.(*Redis.RedisConnect), model)
	} else {
		fmt.Println("query", model.String())
	}
	// model.Set("world")
	// RedisString.Update(conn, model)
}
