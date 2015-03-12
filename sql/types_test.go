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

	DeviceLocation struct {
		DeviceId int64 `sql:"primary" sql-extend:"Device"`
		Location string
	}

	Object struct {
		Id       int64 `sql:"primary,autoincrement"`
		DeviceId int64 `sql-child:"Device"`
		Name     string
	}

	Indicator struct {
		Id       int64 `sql:"primary,autoincrement"`
		ObjectId int64 `sql-child:"Object"`
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

/*
This demonstrates how to update a record using the API.  You will notice that the API requires objects to be passed by reference frequently.  This is to allow the API to be more flexible, which we'll see in the next section
*/
func Example_3() {
	//Create a database with blank tables, for example only
	service := tableCreationBoilerplate()

	//Insert a device into the database
	service.Create(Device{Name: "Device 1"})

	device := Device{}
	//Create an iterator over all of the devices
	next, _ := service.ReadAll(device)
	fmt.Println("Device at the start")
	//Iterate over the deivces
	for next(&device) {
		fmt.Println(device)
	}

	//Passing the record in by reference is very important.  It allows the API to be more flexible, which we'll see later
	service.ReadById(1, &device)

	device.Name = "A New Name"
	service.Update(&device)

	next, _ = service.ReadAll(device)
	fmt.Println("Device at the end")
	//Iterate over the deivces
	for next(&device) {
		fmt.Println(device)
	}

	//Output:
	//Device at the start
	//{1 Device 1}
	//Device at the end
	//{1 A New Name}
}

/*
The API is polymorphic, so it is possible to use the same service to query multiple types of records.  This is why it is important to always pass a pointer to the API
*/
func Example_4() {
	//Create a database with blank tables, for example only
	service := tableCreationBoilerplate()

	//Insert a device into the database
	service.Create(Device{Name: "Device 1"})

	//And now we insert an object into the database
	service.Create(Object{Name: "Object 1", DeviceId: 1})

	device := Device{}
	//Create an iterator over all of the devices
	next, _ := service.ReadAll(device)
	fmt.Println("Devices")
	//Iterate over the deivces
	for next(&device) {
		fmt.Println(device)
	}

	object := Object{}
	//Create an iterator over all of the objects
	next, _ = service.ReadAll(object)
	fmt.Println("Objects")
	//Iterate over the deivces
	for next(&object) {
		fmt.Println(object)
	}

	//Output:
	//Devices
	//{1 Device 1}
	//Objects
	//{1 1 Object 1}
}
