package hash

import (
	"reflect"
	"strconv"
)

func GetValue(value reflect.Value) string {
	switch reflect.TypeOf(reflect.Indirect(value).Interface()).Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.String:
		return value.String()
	// case reflect.Slice:
	// 	return Lib.ToString(value)
	case reflect.Bool:
		if value.Bool() {
			return "1"
		} else {
			return "0"
		}
	default:
		return ""
	}

}

type HashObject struct {
	key     string
	class   reflect.Type
	data    reflect.Value
	updates map[string]interface{}
	// id      int64 //计数 识别对象用的
}

func (hash *HashObject) Set(key string, value interface{}) error {
	hash.updates[key] = value
	return nil
}
func (hash *HashObject) Get(key string) interface{} {
	if hash.updates[key] != nil {
		return hash.updates[key]
	}
	return hash.data.FieldByName(key).Interface()
}
func (hash *HashObject) UpdateString() string {
	str := " HMSet " + hash.key
	for key, value := range hash.updates {
		field, b := hash.class.FieldByName(key)
		if !b {
			continue
		}
		name := field.Tag.Get("name")
		if name == "" {
			name = field.Name
		}
		str += " " + name + " " + GetValue(reflect.ValueOf(value))
	}
	return str
}
func (hash *HashObject) AddString() string {
	str := "HMSET " + hash.key + " "
	size := hash.class.NumField()
	for i := 0; i < size; i++ {
		field := hash.class.Field(i)
		name := field.Tag.Get("name")
		if name == "" {
			name = field.Name
		}
		v := hash.data.Elem().Field(i)
		str += name + " " + GetValue(v) + " "
	}
	return str
}
