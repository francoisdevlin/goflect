package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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

func TestBasicTableOpsBar(t *testing.T) {
	c, _ := sql.Open("sqlite3", ":memory:")
	service := SqliteRecordService{c}
	message := CreateSQLiteTable(&Bar{})
	_, err := c.Exec(message)
	if err != nil {
		t.Error("Miss creating table")
	}

	mocker := MockerStruct{SkipId: true}
	service.Insert((mocker.Mock(1, &Bar{})))
	service.Insert((mocker.Mock(2, &Bar{})))
	service.Insert((mocker.Mock(3, &Bar{})))
	service.Insert((mocker.Mock(4, &Bar{})))

	mocker = MockerStruct{SkipId: false}
	retrieved := Bar{}
	expected := Bar{}
	service.Get(1, &retrieved)
	mocker.Mock(1, &expected)

	if retrieved != expected {
		t.Error("Error on first record equality")
	}

	next := service.ReadAll(&Bar{})
	i := 0
	for next(&retrieved) {
		i++
		mocker.Mock(int64(i), &expected)
		if retrieved != expected {
			t.Error(fmt.Sprintf("Error with autoincrement, Id: %v B: %v", retrieved.Id, retrieved.B))
		}
	}

}
