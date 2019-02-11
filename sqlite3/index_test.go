package sqlite3

import (
	"fmt"
	"testing"
)

type UserInfo struct {
	Uid      int64  `sql:"pk" name:"id" tname:"userInfo" pk:"1"`
	Username string `sql:"nn" name:"username" `
}

func Test_CreateTable(t *testing.T) {
	pool, err := CreateSqlite3Pool(map[string]interface{}{
		"min": 1,
		"max": 2,
		"db":  "./test.db",
	})
	if err != nil {
		t.Error(err)
	}
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
	}
	if conn == nil {
		t.Error("conn 有问题")
	}
	c := conn.(*Sqlite3Connect)
	model := CreateModel((*UserInfo)(nil))
	err = model.CreateTable(c)
	if err != nil {
		t.Error(err)
	}
	// defer pool.Release()
}

func Test_Add(t *testing.T) {
	pool, err := CreateSqlite3Pool(map[string]interface{}{
		"min": 1,
		"max": 2,
		"db":  "./test.db",
	})
	if err != nil {
		t.Error(err)
	}
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
	}
	connect := conn.(*Sqlite3Connect)
	model := CreateModel((*UserInfo)(nil))
	err = model.Add(connect, &UserInfo{
		Username: "tubaoge",
	})
	if err != nil {
		t.Error(err)
	}
}

// func Test_Query(t *testing.T) {
// 	pool, err := CreateSqlite3Pool(map[string]interface{}{
// 		"min": 1,
// 		"max": 2,
// 		"db":  "./test.db",
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	conn, err := pool.Get()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	connect := conn.(*Sqlite3Connect)
// 	model := CreateModel((*UserInfo)(nil))
// 	datas, err := model.Query(connect, map[string]QueryObject{
// 		"username": QueryObject{
// 			Ne: "tubaoge",
// 		},
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// useInfos, err := datas.([]UserInfo)
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }
// 	// useInfos := make([]UserInfo, len(datas))
// 	// for i, arg := range datas {
// 	// 	useInfos[i] = arg.(UserInfo)
// 	// }

// 	fmt.Println(datas[0].Interface().(UserInfo).Uid)

// }

// func Test_Delete(t *testing.T) {
// 	pool, err := CreateSqlite3Pool(map[string]interface{}{
// 		"min": 1,
// 		"max": 2,
// 		"db":  "./test.db",
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	conn, err := pool.Get()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	connect := conn.(*Sqlite3Connect)
// 	model := CreateModel((*UserInfo)(nil))
// 	err = model.Delete(connect, &UserInfo{
// 		Uid: 1,
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func Test_Update(t *testing.T) {
	pool, err := CreateSqlite3Pool(map[string]interface{}{
		"min": 1,
		"max": 2,
		"db":  "./test.db",
	})
	if err != nil {
		t.Error(err)
	}
	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
	}
	connect := conn.(*Sqlite3Connect)
	model := CreateModel((*UserInfo)(nil))
	datas, err := model.Query(connect, map[string]QueryObject{
		"username": QueryObject{
			Ne: "tubaoge",
		},
	})
	if err != nil {
		t.Error(err)
	}
	if len(datas) == 0 {
		return
	}
	datas[0].Set("Username", "shabi")
	fmt.Println(datas[0])
	err = model.Update(connect, datas[0])
	if err != nil {
		t.Error(err)
	}
}
