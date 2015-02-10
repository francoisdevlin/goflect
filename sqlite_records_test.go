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
	mocker = MockerStruct{SkipId: true}
	service.Get(1, &temp)
	mocker.Mock(10, &temp)
	service.Update(&temp)
	service.Get(1, &temp)
	if (temp != Foo{Id: 1, A: "10th", B: 10}) {
		t.Error("Update Failed")
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

	MAX_COUNT := 4 //Why not 4?
	mocker := MockerStruct{SkipId: true}
	for i := 0; i < MAX_COUNT; i++ {
		service.Insert((mocker.Mock(int64(i+1), retrieved)))
	}

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
	if i != MAX_COUNT {
		t.Errorf("Too few records found, expected %v, found %v", MAX_COUNT, i)
	}

	mocker.Mock(1, expected)
	mocker = MockerStruct{SkipId: true}
	service.Get(1, retrieved)
	mocker.Mock(10, retrieved)
	mocker.Mock(10, expected)
	service.Update(retrieved)
	service.Get(1, retrieved)
	if !reflect.DeepEqual(retrieved, expected) {
		t.Error("Error on first record update")
	}
	service.Delete(retrieved)

	next = service.ReadAll(retrieved)
	i = 0
	for next(retrieved) {
		i++
	}
	if i != MAX_COUNT-1 {
		t.Errorf("Too few records found, expected %v, found %v", MAX_COUNT-1, i)
	}

	service.DeleteById(2, retrieved)

	next = service.ReadAll(retrieved)
	i = 0
	for next(retrieved) {
		i++
	}
	if i != MAX_COUNT-2 {
		t.Errorf("Too few records found, expected %v, found %v", MAX_COUNT-2, i)
	}

}

func TestBasicTableOpsBar(t *testing.T) {
	basicWriteHelper(t, &Bar{}, &Bar{})
}

func TestBasicTableOpsFoo2(t *testing.T) {
	basicWriteHelper(t, &Foo{}, &Foo{})
}

func TestBasicTableOpsInts(t *testing.T) {
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

func TestBasicTableOpsUints(t *testing.T) {
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

func TestBasicTableOpsFloats(t *testing.T) {
	type Baz struct {
		Id  int64 `sql:"primary,autoincrement"`
		F32 float32
		F64 float64
	}
	basicWriteHelper(t, &Baz{}, &Baz{})
}
