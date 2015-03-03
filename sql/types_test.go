package records

import (
	"database/sql"
	"fmt"
	"reflect"
)

//Use our types to define a schema
type (
	Device struct {
		Id   int64 `sql:"primary,autoincrement"`
		Name string
	}
	Object struct {
		Id       int64 `sql:"primary,autoincrement"`
		DeviceId int64
		Name     string
	}
)

func tableCreationBoilerplate() RecordService {
	c, _ := sql.Open("sqlite3", ":memory:")
	service := NewSqliteService(c)
	sqlService, _ := service.delegate.(Definer)
	sqlService.Define(&Device{})
	sqlService.Define(&Object{})
	return service
}

func printAll(service RecordService, record interface{}) {
	typ := reflect.TypeOf(record)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	next, _ := service.ReadAll(record)
	fmt.Println(typ.Name() + "s")
	for next(record) {
		val := reflect.ValueOf(record)
		fmt.Println(val.Elem().Interface())
	}
}

/*
There are many items that we will be using throughout the examples.  This is a reference of the functions
*/
func Example_reference() {
	type (
		Device struct {
			Id   int64 `sql:"primary,autoincrement"`
			Name string
		}
		Object struct {
			Id       int64 `sql:"primary,autoincrement"`
			DeviceId int64
			Name     string
		}
	)

	//This function will be used at the beginning of many examples to
	tableCreationBoilerplate := func() RecordService {
		c, _ := sql.Open("sqlite3", ":memory:")
		service := NewSqliteService(c)
		sqlService, _ := service.delegate.(Definer)
		sqlService.Define(&Device{})
		sqlService.Define(&Object{})
		return service
	}

	printAll := func(service RecordService, record interface{}) {
		typ := reflect.TypeOf(record)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		next, _ := service.ReadAll(record)
		fmt.Println(typ.Name() + "s")
		for next(record) {
			val := reflect.ValueOf(record)
			fmt.Println(val.Elem().Interface())
		}
	}

	//A quick example of using the functions to print records
	service := tableCreationBoilerplate()
	service.Create(Device{Name: "Device 1"})

	device := Device{}
	printAll(service, &device)

	//Output:
	//Devices
	//{1 Device 1}
}
