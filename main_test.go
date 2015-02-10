package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
)

type Foo struct {
	Id int64  `sql:"primary,autoincrement"`
	A  string `sql:"unique,nominal"`
	B  int64  `desc:"This is a human readable description"`
}

type Bar struct {
	Id uint64 `sql:"primary,autoincrement"`
	A  string `sql:"unique,nominal"`
	B  bool   `desc:"This is a human readable description"`
}

/*****
 * Begin The Tests
  ****/
func TestSqliteTableCreate(t *testing.T) {
	c, _ := sql.Open("sqlite3", ":memory:")
	message := CreateSQLiteTable(&Foo{})
	_, err := c.Exec(message)
	if err != nil {
		t.Error("Miss creating table")
	}
	_, err = c.Exec(message)
	if err != nil {
		t.Error("Miss recreating creating table")
	}
}

func TestBasicTableOpsFoo(t *testing.T) {
	c, _ := sql.Open("sqlite3", ":memory:")
	service := SqliteRecordService{c}
	message := CreateSQLiteTable(&Foo{})
	_, err := c.Exec(message)
	if err != nil {
		t.Error("Miss creating table")
	}

	mocker := MockerStruct{SkipId: true}
	service.Insert((mocker.Mock(1, &Foo{})))
	service.Insert((mocker.Mock(2, &Foo{})))
	service.Insert((mocker.Mock(3, &Foo{})))
	service.Insert((mocker.Mock(4, &Foo{})))

	temp := Foo{}
	service.Get(1, &temp)

	if (temp != Foo{Id: 1, A: "1st", B: 1}) {
		t.Error("Error Retrieving Record")
	}

	next := service.ReadAll(&Foo{})
	for next(&temp) {
		if temp.Id != temp.B {
			t.Error(fmt.Sprintf("Error with autoincrement, Id: %v B: %v", temp.Id, temp.B))
		}
	}

	next = service.ReadAllWhere(&Foo{}, map[string]interface{}{
		"B": 1,
	})

	for next(&temp) {
		if temp.A != "1st" {
			t.Error(fmt.Sprintf("Error with ID lookup, Id: %v A: %v", temp.Id, temp.A))
		}
	}

	service.GetByNominal("2nd", &temp)

	for next(&temp) {
		if temp.B != 2 {
			t.Error(fmt.Sprintf("Error with lookup, Id: %v A: %v", temp.Id, temp.A))
		}
	}
}

func basicWriteHelper(t *testing.T, retrieved, expected interface{}) {
	c, _ := sql.Open("sqlite3", ":memory:")
	service := SqliteRecordService{c}
	message := CreateSQLiteTable(retrieved)
	_, err := c.Exec(message)
	if err != nil {
		t.Error("Miss creating table")
	}

	mocker := MockerStruct{SkipId: true}
	service.Insert((mocker.Mock(1, retrieved)))
	service.Insert((mocker.Mock(2, retrieved)))
	service.Insert((mocker.Mock(3, retrieved)))
	service.Insert((mocker.Mock(4, retrieved)))

	mocker = MockerStruct{SkipId: false}
	service.Get(1, retrieved)
	mocker.Mock(1, expected)

	if !reflect.DeepEqual(retrieved, expected) {
		t.Error("Error on first record equality")
	}

	next := service.ReadAll(retrieved)
	i := 0
	for next(retrieved) {
		i++
		mocker.Mock(int64(i), expected)
		if !reflect.DeepEqual(retrieved, expected) {
			t.Error(fmt.Sprintf("Error with autoincrement, R: %v E: %v", retrieved, expected))
		}
	}

}

func TestBasicTableOpsBar(t *testing.T) {
	basicWriteHelper(t, &Bar{}, &Bar{})
}

func TestBasicTableOpsFoo2(t *testing.T) {
	basicWriteHelper(t, &Foo{}, &Foo{})
}

func TestBasicTableOpsFoo3(t *testing.T) {
	type Baz struct {
		Id  int64 `sql:"primary,autoincrement"`
		I   int
		I64 int64
		I32 int32
		I16 int16
		I8  int8
	}
	basicWriteHelper(t, &Baz{}, &Baz{})
}

func TestBasicTableOpsFoo4(t *testing.T) {
	type Baz struct {
		Id  int64 `sql:"primary,autoincrement"`
		U   uint
		U64 uint64
		U32 uint32
		U16 uint16
		U8  uint8
	}
	basicWriteHelper(t, &Baz{}, &Baz{})
}
