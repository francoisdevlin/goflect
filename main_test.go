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
	B  string `desc:"This is a human readable description"`
}

/*****
 * Begin The Tests
  ****/
func TestSqliteTableCreate(t *testing.T) {
	fmt.Println("Hello Test")
	c, _ := sql.Open("sqlite3", ":memory:")
	message := CreateSQLiteTable(&Foo{})
	_, err := c.Exec(message)
	if err != nil {
		fmt.Println("Miss creating table")
	}
	message = InsertSQLiteRecord(&Foo{A: "First", B: "Second"})
	_, err = c.Exec(message)
	message = InsertSQLiteRecord(&Foo{A: "First", B: "Second"})
	_, err = c.Exec(message)
	message = InsertSQLiteRecord(&Foo{A: "First", B: "Second"})
	_, err = c.Exec(message)
	message = InsertSQLiteRecord(&Foo{A: "First", B: "Second"})
	_, err = c.Exec(message)
	message = ListSQLiteRecord(&Foo{})
	//fmt.Println(message)
	rows, err := c.Query(message)

	if err != nil {
		fmt.Println("Query Error", err)
	} else {
		temp := Foo{}
		for NextRow(rows, &temp) {
			fmt.Println(temp)
		}
	}

	message = ListSQLiteNominal(&Foo{})
	rows, err = c.Query(message)

	if err != nil {
		fmt.Println("Query Error", err)
	} else {
		temp := Nominal{}
		for NextRow(rows, &temp) {
			fmt.Println(temp)
		}
	}

}
