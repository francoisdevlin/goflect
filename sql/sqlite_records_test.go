package records

import (
	"database/sql"
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"git.sevone.com/sdevlin/goflect.git/mock"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
)

/*
This is a basic example showing how the metadata maps to a sqlite service table creation.  It needs to be fixed to allow unique keys to be respected
*/
func ExampleSqlDefiner_sqliteBasic() {
	type Foo struct {
		Id int64  `sql:"primary,autoincrement"`
		A  string `sql:"unique,nominal"`
		B  int64
	}

	c, _ := sql.Open("sqlite3", ":memory:")
	service := NewSqliteService(c)
	sqlService, _ := service.delegate.(SqlDefiner)
	fmt.Println(sqlService.CreateStatement(Foo{}))

	//Output:
	//CREATE TABLE IF NOT EXISTS Foo(
	//	`Id` integer primary key autoincrement not null,
	//	`A` string not null,
	//	`B` integer
	//)
}

/*
This is a verbose example of how to use the Create API, with an accompanying read section for verification
*/
func ExampleRecordService_createSqliteVerbose() {
	type Foo struct {
		Id int64  `sql:"primary,autoincrement"`
		A  string `sql:"unique,nominal"`
		B  int64
	}

	c, _ := sql.Open("sqlite3", ":memory:")
	service := NewSqliteService(c)
	sqlService, _ := service.delegate.(Definer)
	err := sqlService.Define(&Foo{})
	if err == nil {
		fmt.Println("Table created properly")
	}

	foo := Foo{A: "Hello World", B: 10}
	//Passing a pointer is important, so that the id can be placed in foo
	err = service.Create(&foo)
	fmt.Println(foo)
	if err == nil {
		fmt.Println("Record Createed properly")
	}
	next, err := service.ReadAll(&foo)
	if err == nil {
		fmt.Println("Records read properly")
	}
	newFoo := Foo{}
	//The next function is an iterator, which stores the result in newFoo
	//It returns true if there are more records to process
	//It will close the statement when it is complete
	//It must take a pointer
	for next(&newFoo) {
		fmt.Println(newFoo)
	}

	//Output:
	//Table created properly
	//{0 Hello World 10}
	//Record Createed properly
	//Records read properly
	//{1 Hello World 10}
}

/*
This is a more idiomatic example of using the create function
*/
func ExampleRecordService_createSqliteIdiomatic() {
	type Foo struct {
		Id int64  `sql:"primary,autoincrement"`
		A  string `sql:"unique,nominal"`
		B  int64
	}

	c, _ := sql.Open("sqlite3", ":memory:")
	service := NewSqliteService(c)
	sqlService, _ := service.delegate.(Definer)
	err := sqlService.Define(&Foo{})
	if err != nil {
		return
	}

	foo := Foo{A: "Hello World", B: 10}
	err = service.Create(&foo)
	if err != nil {
		return
	}
	next, err := service.ReadAll(&foo)
	if err != nil {
		return
	}
	newFoo := Foo{}
	for next(&newFoo) {
		fmt.Println(newFoo)
	}

	//Output:
	//{1 Hello World 10}
}

/*
This is a more idiomatic example of using the create function
*/
func ExampleRecordService_infoTest() {

	c, _ := sql.Open("sqlite3", ":memory:")
	service := NewSqliteService(c)
	sqlService, _ := service.delegate.(Definer)
	err := sqlService.Define(&goflect.Info{})
	if err != nil {
		fmt.Println(err)
		return
	}

	//Output:
	//{1 Hello World 10}
}

/*
This shows the power of the rails convention transform in action
*/
func Example_railsConvention1() {
	//See Package reference example for definition
	service := tableCreationBoilerplate()

	//create our first device
	service.Create(Device{Name: "Device 1"})

	//Print all the devices
	device := Device{}
	printAll(service, &device)

	//Create a new data service with the RailsConvention transform
	deviceService := NewTransformService(RailsConvention(device), service)

	//Create some objects
	deviceService.Create(&Object{Name: "Object 1"})
	deviceService.Create(&Object{Name: "Object 2"})
	deviceService.Create(&Object{Name: "Object 3"})

	//And notice that the device id was handled for us automatically
	object := Object{}
	printAll(service, &object)

	//Output:
	//Devices
	//{1 Device 1}
	//Objects
	//{1 1 Object 1}
	//{2 1 Object 2}
	//{3 1 Object 3}
}

/*****
 * Begin The Tests
  ****/
func TestBasicTableOpsFoo(t *testing.T) {
	type Foo struct {
		Id int64  `sql:"primary,autoincrement"`
		A  string `sql:"unique,nominal"`
		B  int64
	}

	c, _ := sql.Open("sqlite3", ":memory:")
	service := NewSqliteService(c)
	sqlService, _ := service.delegate.(Definer)
	err := sqlService.Define(&Foo{})
	if err != nil {
		t.Error("Miss creating table")
	}

	mocker := mock.MockerStruct{SkipId: true}
	service.Create((mocker.Mock(1, &Foo{})))
	service.Create((mocker.Mock(2, &Foo{})))
	service.Create((mocker.Mock(3, &Foo{})))
	service.Create((mocker.Mock(4, &Foo{})))

	temp := Foo{}
	service.Get(1, &temp)

	if (temp != Foo{Id: 1, A: "1st", B: 1}) {
		t.Error("Error Retrieving Record")
	}

	next, _ := service.ReadAll(&Foo{})
	for next(&temp) {
		if temp.Id != temp.B {
			t.Error(fmt.Sprintf("Error with autoincrement, Id: %v B: %v", temp.Id, temp.B))
		}
	}

	match := matcher.NewStructMatcher()
	match.AddField("B", matcher.Eq(1))
	next, _ = service.ReadAllWhere(&Foo{}, match)

	for next(&temp) {
		if temp.A != "1st" {
			t.Error(fmt.Sprintf("Error with ID lookup, Id: %v A: %v", temp.Id, temp.A))
		}
	}

	mocker = mock.MockerStruct{SkipId: true}
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
	service := NewSqliteService(c)
	sqlService, _ := service.delegate.(Definer)
	err := sqlService.Define(retrieved)
	if err != nil {
		t.Error("Miss creating table")
	}

	MAX_COUNT := 4 //Why not 4?
	mocker := mock.MockerStruct{SkipId: true}
	for i := 0; i < MAX_COUNT; i++ {
		service.Create((mocker.Mock(int64(i+1), retrieved)))
	}

	mocker = mock.MockerStruct{SkipId: false}
	service.Get(1, retrieved)
	mocker.Mock(1, expected)
	//fmt.Println(retrieved)
	//fmt.Println(expected)

	if !reflect.DeepEqual(retrieved, expected) {
		t.Error("Error on first record equality")
	}

	next, _ := service.ReadAll(retrieved)
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
	mocker = mock.MockerStruct{SkipId: true}
	service.Get(1, retrieved)
	mocker.Mock(10, retrieved)
	mocker.Mock(10, expected)
	service.Update(retrieved)
	service.Get(1, retrieved)
	if !reflect.DeepEqual(retrieved, expected) {
		t.Error("Error on first record update")
	}
	service.Delete(retrieved)

	next, _ = service.ReadAll(retrieved)
	i = 0
	for next(retrieved) {
		i++
	}
	if i != MAX_COUNT-1 {
		t.Errorf("Too few records found, expected %v, found %v", MAX_COUNT-1, i)
	}

	service.DeleteById(2, retrieved)

	next, _ = service.ReadAll(retrieved)
	i = 0
	for next(retrieved) {
		i++
	}
	if i != MAX_COUNT-2 {
		t.Errorf("Too few records found, expected %v, found %v", MAX_COUNT-2, i)
	}

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

func TestBasicTableOpsBasicEmbed(t *testing.T) {
	type Embed struct {
		F32 float32
		F64 float64
	}
	type Baz struct {
		Id int64 `sql:"primary,autoincrement"`
		Embed
	}
	basicWriteHelper(t, &Baz{}, &Baz{})
}

func TestSQLDeepEmbed(t *testing.T) {
	type E00 struct{ I00 int64 }
	type E01 struct{ E00, I01 int64 }
	type E02 struct{ E01, I02 int64 }
	type E03 struct{ E02, I03 int64 }
	type E04 struct{ E03, I04 int64 }
	type E05 struct{ E04, I05 int64 }
	type E06 struct{ E05, I06 int64 }
	type E07 struct{ E06, I07 int64 }
	type E08 struct{ E07, I08 int64 }
	type E09 struct{ E08, I09 int64 }
	type E10 struct{ E09, I10 int64 }
	type E11 struct{ E10, I11 int64 }
	type E12 struct{ E11, I12 int64 }
	type E13 struct{ E12, I13 int64 }
	type E14 struct{ E13, I14 int64 }
	type E15 struct{ E14, I15 int64 }
	type E16 struct{ E15, I16 int64 }
	type E17 struct{ E16, I17 int64 }
	type E18 struct{ E17, I18 int64 }
	type E19 struct{ E18, I19 int64 }
	type E20 struct{ E19, I20 int64 }
	type Baz struct {
		Id int64 `sql:"primary,autoincrement"`
		E20
	}
	basicWriteHelper(t, &Baz{}, &Baz{})
}

func TestSQLEmbedAttributes(t *testing.T) {
	type IdStruct struct {
		Id int64 `sql:"primary,autoincrement"`
	}
	type BazStruct struct {
		F32 float32
		F64 float64
	}
	type Baz struct {
		IdStruct
		BazStruct
	}
	basicWriteHelper(t, &Baz{}, &Baz{})
}
