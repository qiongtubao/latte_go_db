package redis

import (
	"errors"
	"strings"

	"github.com/go-redis/redis"
)

type RedisConnect struct {
	Client *redis.Client
}

func (connect *RedisConnect) Close() {
	connect.Close()
}
func (connect *RedisConnect) Get(key string) (interface{}, error) {
	return connect.Client.Get(key).Result()
}
func (connect *RedisConnect) Set(key string, value interface{}) (interface{}, error) {
	return nil, connect.Client.Set(key, value, 0).Err()
}

func (connect *RedisConnect) HMSet(key string, value []string) (interface{}, error) {
	// m := map[string]interface{}{}
	kv_map := make(map[string]interface{})
	size := len(value)
	for i := 0; i < size; i += 2 {
		kv_map[value[i]] = value[i+1]
	}
	return nil, connect.Client.HMSet(key, kv_map).Err()
}
func (connect *RedisConnect) HMGet(key string, value []string) (interface{}, error) {
	return connect.Client.HMGet(key, value...).Result()
}
func (connect *RedisConnect) Exec(str string) (interface{}, error) {
	strs := strings.Fields(str)
	switch strings.ToUpper(strs[0]) {
	case "GET":
		if len(strs) < 2 {
			return nil, errors.New("Missing parameters ")
		}
		return connect.Get(strs[1])
	case "SET":
		if len(strs) < 3 {
			return nil, errors.New("Missing parameters ")
		}
		return connect.Set(strs[1], strs[2])
	case "HMSET":
		if len(strs) < 4 {
			return nil, errors.New("Missing parameters ")
		}
		if len(strs)%2 != 0 {
			return nil, errors.New(" parameters error")
		}
		return connect.HMSet(strs[1], strs[2:])
	case "HMGET":
		if len(strs) < 3 {
			return nil, errors.New("Missing parameters ")
		}
		return connect.HMGet(strs[1], strs[2:])
	}
	return nil, errors.New("命令错误")
}
