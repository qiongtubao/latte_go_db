package string

import (
	Redis "github.com/qiongtubao/latte_go_db/redis"
)

type Model struct {
	key   string
	value interface{}
}

func (model *Model) Set(value interface{}) {
	model.value = value
}
func (model *Model) Interface() interface{} {
	return model.value
}
func (model *Model) String() string {
	if model.value == nil {
		return ""
	}
	return model.value.(string)
}

func CreateModel(k string) *Model {
	return &Model{
		key: k,
	}
}

func Get(conn *Redis.RedisConnect, k string) (*Model, error) {
	value, err := conn.Get(k)
	if err != nil {
		return nil, err
	}
	return &Model{
		key:   k,
		value: value,
	}, nil
}

func Update(conn *Redis.RedisConnect, model *Model) error {
	_, err := conn.Set(model.key, model.value)
	return err
}
