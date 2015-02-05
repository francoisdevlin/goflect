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

func TestBasicTableOps(t *testing.T) {
	c, _ := sql.Open("sqlite3", ":memory:")
	service := SqliteRecordService{c}
	message := CreateSQLiteTable(&Foo{})
	_, err := c.Exec(message)
	if err != nil {
		t.Error("Miss creating table")
	}
	service.Insert(&Foo{A: "First", B: 1})
	service.Insert(&Foo{A: "First", B: 1})
	service.Insert(&Foo{A: "First", B: 2})
	service.Insert(&Foo{A: "Bacon", B: 1})

	next := service.ReadAll(&Foo{})

	temp := Foo{}
	for next(&temp) {
		fmt.Println(temp)
	}

	next = service.ReadAllWhere(&Foo{}, map[string]interface{}{
		"B": 1,
	})
	for next(&temp) {
		fmt.Println(temp)
	}

	service.Get(1, &temp)
	fmt.Println(temp)

	service.GetByNominal("Bacon", &temp)
	fmt.Println(temp)

	nextNom := service.ReadAllNominal(&Foo{})
	tempNom := Nominal{}
	for nextNom(&tempNom) {
		fmt.Println(tempNom)
	}
}
