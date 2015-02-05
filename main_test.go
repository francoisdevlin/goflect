package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func prepConn() *sql.DB {
	c, _ := sql.Open("sqlite3", ":memory:")
	return c
}

type Foo struct {
	Id int64  `sql:"primary,autoincrement"`
	A  string `sql:"unique,nominal"`
	B  int64  `desc:"This is a human readable description"`
}

/*****
 * Begin The Tests
  ****/
func TestSqliteTableCreate(t *testing.T) {
	fmt.Println("Hello Test")
	c, _ := sql.Open("sqlite3", ":memory:")
	service := SqliteRecordService{c}
	message := CreateSQLiteTable(&Foo{})
	_, err := c.Exec(message)
	if err != nil {
		fmt.Println("Miss creating table")
	}
	service.Insert(&Foo{A: "First", B: 1})
	service.Insert(&Foo{A: "First", B: 1})
	service.Insert(&Foo{A: "First", B: 1})
	service.Insert(&Foo{A: "First", B: 1})

	next := service.ReadAll(&Foo{})

	temp := Foo{}
	for next(&temp) {
		fmt.Println(temp)
	}

	nextNom := service.ReadAllNominal(&Foo{})
	tempNom := Nominal{}
	for nextNom(&tempNom) {
		fmt.Println(tempNom)
	}

}
