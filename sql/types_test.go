package records

import (
	"database/sql"
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
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

/*
Let's see how easy it is to insert records into a blank database
*/
func Example_1() {
	//Create a database with blank tables, for example only
	service := tableCreationBoilerplate()

	//Insert a device into the database
	service.Create(Device{Name: "Device 1"})

	device := Device{}
	//Create an iterator over all of the devices
	next, _ := service.ReadAll(device)
	fmt.Println("Devices")
	//Iterate over the deivces
	for next(&device) {
		fmt.Println(device)
	}

	//Output:
	//Devices
	//{1 Device 1}
}

/*
Let's insert a few devices into the database, and the perform a query on the database
*/
func Example_2() {
	//Create a database with blank tables, for example only
	service := tableCreationBoilerplate()

	//Insert a device into the database
	service.Create(Device{Name: "Device 1"})
	service.Create(Device{Name: "Device 2"})
	service.Create(Device{Name: "Device 3"})

	device := Device{}
	//Create an iterator over all of the devices
	next, _ := service.ReadAll(device)
	fmt.Println("Devices")
	//Iterate over the deivces
	for next(&device) {
		fmt.Println(device)
	}

	//Let's find device 1
	//We'll need to create a matcher the describes the set that this device is in
	query := matcher.NewStructMatcher()
	query.AddField("Name", matcher.Eq("Device 1"))
	next, _ = service.ReadAllWhere(device, query)

	fmt.Println("Our device")
	//Iterate over the deivces
	for next(&device) {
		fmt.Println(device)
	}

	//Finding a device that doesn't exist will have the iterator quit immediately
	query = matcher.NewStructMatcher()
	query.AddField("Name", matcher.Eq("DOES NOT EXIST"))
	next, _ = service.ReadAllWhere(device, query)

	fmt.Println("No Devices Match Query")
	//Iterate over the deivces
	for next(&device) {
		fmt.Println(device)
	}

	//Output:
	//Devices
	//{1 Device 1}
	//{2 Device 2}
	//{3 Device 3}
	//Our device
	//{1 Device 1}
	//No Devices Match Query
}
