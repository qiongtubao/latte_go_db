package hash

import (
	"fmt"
	"testing"

	Redis "github.com/qiongtubao/latte_go_db/redis"
)

type UserInfo struct {
	Name string `name:"name" `
	Age  int64  `name:"age" `
}

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

	model := CreateModel((*UserInfo)(nil))
	_, err = model.Add(conn.(*Redis.RedisConnect), "testHash", &UserInfo{
		Name: "tubaoge",
		Age:  1,
	})
	if err != nil {
		t.Error(err)
	}

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

	model := CreateModel((*UserInfo)(nil))
	obj, err := model.Query(conn.(*Redis.RedisConnect), "testHash")
	if err != nil {
		t.Error(err)
	}
	obj.Set("Age", 11)
	// obj.Set("Age", 11)
	err = model.Update(conn.(*Redis.RedisConnect), obj)
	if err != nil {
		t.Error(err)
	}
	o, err := model.Query(conn.(*Redis.RedisConnect), "testHash")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(o.Get("Age"))
}

// func Test_Del(t *testing.T) {
// 	pool, err := Redis.CreateRedisPool(map[string]interface{}{
// 		"min":  1,
// 		"max":  2,
// 		"host": "localhost",
// 		"port": "6379",
// 		"db":   0,
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	conn, err := pool.Get()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	model := CreateModel((*UserInfo)(nil))
// 	model.Del(conn.(*Redis.RedisConnect), "testHash")
// }

// data, err := model.Query(conn.(*Redis.RedisConnect), "hash")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	data.Set("Age", 10)
// 	model.Update(data)
