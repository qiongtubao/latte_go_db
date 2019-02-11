package sqlite3

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Model struct {
	modeClass reflect.Type
	keys      []string
	tableName string
}

var SQLDEBUG bool = true

func CreateModel(data interface{}) Model {
	modeClass := reflect.TypeOf(data).Elem()
	len := modeClass.NumField()
	var tableName string = ""
	for i := 0; i < len; i++ {
		field := modeClass.Field(i)
		// identity := field.Tag.Get("identity")
		// if identity != "" {

		// }
		if tableName == "" {
			tableName = field.Tag.Get("tname")
		}
	}
	if tableName == "" {
		tableName = modeClass.Name()
	}
	model := Model{
		modeClass: modeClass,
		tableName: tableName,
	}
	return model
}
func (model Model) CreateTable(connect SqlConnect) error {
	str := "CREATE TABLE IF NOT EXISTS " + model.tableName + "("
	len := model.modeClass.NumField()
	for i := 0; i < len; i++ {
		if i != 0 {
			str += ","
		}
		field := model.modeClass.Field(i)
		name := field.Tag.Get("name")
		if name == "" {
			name = field.Name
		}
		str += name + " "
		t := field.Tag.Get("type")
		if t == "" {
			switch field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				str += "INTEGER "
				break
			case reflect.String:
				str += "VARCHAR(64) "
				break
			}
		} else {
			str += t + " "
		}
		identity := field.Tag.Get("identity")
		fmt.Println(identity)
		if identity != "" {
			str += " identity(" + identity + ")"
		}
		pk := field.Tag.Get("pk")
		if pk != "" && (pk == "1" || pk == "true") {
			str += "PRIMARY KEY  "
		}
	}
	str += ")"
	if SQLDEBUG {
		fmt.Printf("创建数据库表sql: %v", str)
	}
	_, e := connect.Exec(str)
	return e
}

/**
gt >
lt <
gte >=
lte <=
*/
type QueryObject struct {
	Gt  int
	Lt  int
	Lte int
	Gte int
	Ne  interface{}
}

func (q QueryObject) ToString() string {
	switch q.Ne.(type) {
	case string:
		return "(%v = '" + q.Ne.(string) + "')"
		break
	case int:
		return "(%v = " + strconv.Itoa(q.Ne.(int)) + ")"
		break
	}
	return ""
}

func (model Model) Query(connect SqlConnect, querys map[string]QueryObject) ([]Object, error) {
	var resultsSlice []Object

	str := "SELECT "
	len := model.modeClass.NumField()
	keys := ""
	m := map[string]string{}
	for i := 0; i < len; i++ {
		if i != 0 {
			keys += ","
		}
		field := model.modeClass.Field(i)
		name := field.Tag.Get("name")
		if name == "" {
			name = field.Name
		}
		m[name] = field.Name
		keys += name
	}
	wheres := ""
	once := true
	for key, value := range querys {
		s := value.ToString()
		if s == "" {
			continue
		}
		if !once {
			wheres += " and "
		}
		once = false
		wheres += fmt.Sprintf(value.ToString(), key)
	}
	// keys += ")"
	str += keys + " FROM " + model.tableName
	if wheres != "" {
		str += " where " + wheres
	}
	fmt.Println(str)
	res, err := connect.Query(str)
	// defer res.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fields, err := res.Columns()
	if err != nil {
		return nil, err
	}
	for res.Next() {
		// result := reflect.New(model.modeClass).Elem()
		result := CreateObject(model.modeClass)
		var scanResultContainers []interface{}
		for i := 0; i < len; i++ {
			var scanResultContainer interface{}
			scanResultContainers = append(scanResultContainers, &scanResultContainer)
		}
		if err := res.Scan(scanResultContainers...); err != nil {
			return nil, err
		}
		for ii, key := range fields {
			rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))

			// if rawValue.Interface() == nil {
			// 	continue
			// }
			t := reflect.TypeOf(rawValue.Interface())
			// v := reflect.ValueOf(rawValue.Interface())
			switch t.Kind() {
			case reflect.Slice:
				if t.Elem().Kind() == reflect.Uint8 {
					st := string(rawValue.Interface().([]byte))
					// fmt.Println(key, st, result)
					// result.FieldByName(m[key]).Set(reflect.ValueOf(st))
					result.Set(m[key], st)
					break
				}
			case reflect.Int64:
				// result.FieldByName(m[key]).Set(rawValue.Elem())
				result.Set(m[key], rawValue.Elem())
				break
			}
			// result.FieldByName(key).Elem().Set(rawValue)
		}
		result.Empty()
		resultsSlice = append(resultsSlice, result)
	}
	return resultsSlice, nil
}

func isNull(value reflect.Value) bool {
	switch reflect.TypeOf(reflect.Indirect(value).Interface()).Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() == 0 {
			return true
		}
		break
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value.Uint() == 0 {
			return true
		}
		break
	case reflect.String:
		if value.String() == "" {
			return true
		}
	}
	return false
}
func (model Model) Add(connect SqlConnect, data interface{}) error {
	modeClass := reflect.TypeOf(data).Elem()
	if model.modeClass != modeClass {
		return errors.New("与model 类型不同")
	}
	//stmt, err := model.Db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	str := "INSERT INTO " + model.tableName
	keyStr := "("
	valueStr := " values ("
	value := reflect.ValueOf(data).Elem()
	len := model.modeClass.NumField()
	var once = false
	for i := 0; i < len; i++ {
		if isNull(value.Field(i)) {
			continue
		}
		if once == true {
			keyStr += ","
			valueStr += ","
		}
		once = true
		field := modeClass.Field(i)
		name := field.Tag.Get("name")
		if name == "" {
			name = field.Name
		}
		keyStr += name
		valueStr += GetValue(value.Field(i))
	}
	keyStr += ")"
	valueStr += ")"
	str += keyStr + valueStr
	fmt.Println(str)
	_, err := connect.Exec(str)
	return err
}

func (model Model) whereSql(where interface{}) string {
	wheresql := ""
	var t reflect.Type
	switch where.(type) {
	case reflect.Type:
		t = where.(reflect.Type)
		break
	default:
		t = reflect.TypeOf(where)
		break
	}
	switch t {
	case model.modeClass:
		value := reflect.ValueOf(where)
		len := model.modeClass.NumField()
		once := false
		for i := 0; i < len; i++ {
			if isNull(value.Field(i)) {
				continue
			}
			if once == true {
				wheresql += " and "
			}
			once = true
			field := model.modeClass.Field(i)
			name := field.Tag.Get("name")
			if name == "" {
				name = field.Name
			}
			wheresql += "(" + name + " = " + GetValue(value.Field(i)) + ")"
		}
		break
	default:
		fmt.Println("???????")
		break
	}
	return wheresql
}

func (model Model) Delete(connect SqlConnect, where interface{}) error {
	str := "delete from " + model.tableName
	wheresql := model.whereSql(where)
	if wheresql != "" {
		str += " where " + wheresql
	}
	fmt.Println(str)
	_, e := connect.Exec(str)
	return e
}

func (model Model) Update(connect SqlConnect, obj Object) error {
	str := "update  " + model.tableName + " set "
	str += obj.UpdateSql()
	wheresql := model.whereSql(obj.Interface())
	if wheresql != "" {
		str += " where " + wheresql
	}
	fmt.Println(str)
	_, e := connect.Exec(str)
	return e
}
