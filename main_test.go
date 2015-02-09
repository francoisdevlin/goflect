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
	service.Insert(&Foo{A: "First", B: 1})
	service.Insert(&Foo{A: "Second", B: 2})
	service.Insert(&Foo{A: "Third", B: 3})
	service.Insert(&Foo{A: "Fourth", B: 4})

	temp := Foo{}
	service.Get(1, &temp)

	if (temp != Foo{Id: 1, A: "First", B: 1}) {
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
		if temp.A != "First" {
			t.Error(fmt.Sprintf("Error with ID lookup, Id: %v A: %v", temp.Id, temp.A))
		}
	}

	service.GetByNominal("Second", &temp)

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
	service.Insert(&Bar{A: "First", B: true})
	service.Insert(&Bar{A: "Second", B: true})
	service.Insert(&Bar{A: "Third", B: false})
	service.Insert(&Bar{A: "Fourth", B: true})

	temp := Bar{}
	service.Get(1, &temp)

	//fmt.Println(temp.Id, temp.A, temp.B)
	if (temp != Bar{Id: 1, A: "First", B: true}) {
		t.Error("Error on first record equality")
	}

	//next := service.ReadAll(&Bar{})
	//for next(&temp) {
	//if temp.Id != temp.B {
	//t.Error(fmt.Sprintf("Error with autoincrement, Id: %v B: %v", temp.Id, temp.B))
	//}
	//}

	//next = service.ReadAllWhere(&Bar{}, map[string]interface{}{
	//"B": 1,
	//})

	//for next(&temp) {
	//if temp.A != "First" {
	//t.Error(fmt.Sprintf("Error with ID lookup, Id: %v A: %v", temp.Id, temp.A))
	//}
	//}

	//service.GetByNominal("Second", &temp)

	//for next(&temp) {
	//if temp.B != 2 {
	//t.Error(fmt.Sprintf("Error with lookup, Id: %v A: %v", temp.Id, temp.A))
	//}
	//}
}
