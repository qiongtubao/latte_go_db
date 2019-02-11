package sqlite3

import (
	"fmt"
	"reflect"
)

type Object struct {
	modelClass reflect.Type
	data       reflect.Value
	updates    map[string]interface{}
	id         int64
}

func (object Object) Get(key string) interface{} {
	value := object.data.FieldByName(key)
	return value.Interface()
}

func (object Object) Set(key string, value interface{}) error {
	// v := object.data.FieldByName(key)
	// if object.updates[key] == "" {
	object.updates[key] = value
	// }
	// switch value.(type) {
	// case reflect.Value:
	// 	// v.Set(value.(reflect.Value))
	// 	break
	// default:
	// 	// v.Set(reflect.ValueOf(value))
	// 	break
	// }
	// v.Set(reflect.ValueOf(value))
	return nil
}
func (object *Object) Empty() {
	for key, value := range object.updates {
		v := object.data.FieldByName(key)
		switch value.(type) {
		case reflect.Value:
			v.Set(value.(reflect.Value))
			break
		default:
			v.Set(reflect.ValueOf(value))
			break
		}
	}
	object.updates = map[string]interface{}{}
	fmt.Println("clean ", object.id)
}
func (object Object) Interface() interface{} {
	return object.data.Interface()
}

func (object Object) UpdateSql() string {
	str := ""
	once := false
	fmt.Println(object.updates, object.id)
	for k, v := range object.updates {
		if v == "" {
			continue
		}
		if once {
			str += " , "
		}
		once = true
		// va := object.data.FieldByName(k)
		va := reflect.ValueOf(v)
		ka, b := object.modelClass.FieldByName(k)
		if b != false {
			name := ka.Tag.Get("name")
			if name == "" {
				name = k
			}
			str += name + " = " + GetValue(va)
		}
	}
	return str
}

var idNum int64 = 0

func CreateObject(model reflect.Type) Object {
	data := reflect.New(model).Elem()
	idNum = idNum + 1
	return Object{
		modelClass: model,
		data:       data,
		updates:    map[string]interface{}{},
		id:         idNum,
	}
}
