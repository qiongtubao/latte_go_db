package sqlite3

import (
	"reflect"
	"strconv"
)

func GetValue(value reflect.Value) string {
	switch reflect.TypeOf(reflect.Indirect(value).Interface()).Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.String:
		return "'" + value.String() + "'"
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
