package hash

import (
	"errors"
	"reflect"

	Redis "github.com/qiongtubao/latte_go_db/redis"
	Lib "github.com/qiongtubao/latte_go_lib"
)

type Model struct {
	modelClass reflect.Type
}

func (model *Model) Update(conn *Redis.RedisConnect, hash *HashObject) error {
	// fmt.Println(hash.UpdateString())
	_, err := conn.Exec(hash.UpdateString())
	return err
}
func (model *Model) Query(conn *Redis.RedisConnect, key string) (*HashObject, error) {
	str := "HMGet " + key
	size := model.modelClass.NumField()
	for i := 0; i < size; i++ {
		field := model.modelClass.Field(i)
		name := field.Tag.Get("name")
		if name == "" {
			name = field.Name
		}
		str += " " + name
	}
	result, err := conn.Exec(str)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("查询错误")
	}
	value := reflect.New(model.modelClass).Elem()
	strArray := result.([]interface{})
	for i := 0; i < size; i++ {
		v := value.Field(i)
		field := model.modelClass.Field(i)
		vl, err := Lib.TypeChange(strArray[i], field.Type.Kind())
		if err != nil {
			continue
		}
		v.Set(reflect.ValueOf(vl))
	}
	return &HashObject{
		key:     key,
		class:   model.modelClass,
		data:    value,
		updates: map[string]interface{}{},
	}, nil
}

func (model *Model) Add(conn *Redis.RedisConnect, key string, data interface{}) (*HashObject, error) {
	if reflect.TypeOf(data).Elem() != model.modelClass {
		return nil, errors.New("类型不匹配")
	}
	modelObjcet := &HashObject{
		key:     key,
		class:   model.modelClass,
		data:    reflect.ValueOf(data),
		updates: map[string]interface{}{},
	}
	_, err := conn.Exec(modelObjcet.AddString())
	if err != nil {
		return nil, err
	}
	return modelObjcet, nil
}

func CreateModel(data interface{}) *Model {
	modeClass := reflect.TypeOf(data).Elem()
	return &Model{
		modeClass,
	}
}
