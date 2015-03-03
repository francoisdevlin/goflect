package records

import (
	"fmt"
	"testing"
)

/*
This shows how to use the a service to resetrict behavior
*/
func ExampleRecordService_basicDispatch() {
	//A basic printer
	printStats := func(d1, d2 RecordService) {
		dum1, _ := d1.delegate.(*dummyService)
		dum2, _ := d2.delegate.(*dummyService)
		fmt.Println("Creates - dummy 1, dummy 2:", dum1.Creates, dum2.Creates)
	}

	//Create our dummy services
	dummy1 := NewDummyService()
	dummy2 := NewDummyService()

	//A simple dispatch function
	disp := func(record interface{}) (int, error) {
		dummy, _ := record.(dummyService)
		return dummy.Creates % 2, nil
	}

	//Create the dispatch service
	service := NewDispatchService(disp, []RecordService{dummy1, dummy2})

	//We start with nothing
	printStats(dummy1, dummy2)

	//The create only goes to one dummy
	service.Create(dummyService{Creates: 1})
	printStats(dummy1, dummy2)

	//The create goes to the other dummy
	service.Create(dummyService{Creates: 2})
	printStats(dummy1, dummy2)
	//Output:
	//Creates - dummy 1, dummy 2: 0 0
	//Creates - dummy 1, dummy 2: 0 1
	//Creates - dummy 1, dummy 2: 1 1
}

/*
When one of the underlying services returns an error, it is propogated up.  The other services are unaffected
*/
func ExampleRecordService_buggyService() {

	//Create our dummy services
	dummy := NewDummyService()
	buggy := NewBuggyService()

	//A simple dispatch function
	disp := func(record interface{}) (int, error) {
		dummy, _ := record.(dummyService)
		return dummy.Creates % 2, nil
	}

	//Create the dispatch service
	service := NewDispatchService(disp, []RecordService{dummy, buggy})

	//The create goes to the dummy service
	err := service.Create(dummyService{Creates: 0})
	fmt.Println("No error ", err)

	//The create goes to the buggy service, and an error is returned
	err = service.Create(dummyService{Creates: 1})
	fmt.Println(err)
	//Output:
	//No error  <nil>
	//Intentional Create Error
}

/*
When the function returns an error, it is propogated up.  There is also an error if the index is out of range
*/
func ExampleRecordService_buggyDispatch() {

	//Create our dummy services
	dummy := NewDummyService()
	buggy := NewBuggyService()

	//A simple dispatch function
	buggyDispatch := func(record interface{}) (int, error) {
		return 0, RecordError("Dispatch has a bug")
	}

	//Create the dispatch service
	service := NewDispatchService(buggyDispatch, []RecordService{dummy, buggy})

	//The create goes to the dummy service
	err := service.Create(dummyService{Creates: 0})
	fmt.Println(err)

	//A simple dispatch function
	outOfRange := func(record interface{}) (int, error) {
		return 2, nil
	}

	service = NewDispatchService(outOfRange, []RecordService{dummy, buggy})
	err = service.Create(dummyService{Creates: 1})
	fmt.Println(err)
	//Output:
	//Dispatch has a bug
	//Dispatch index out of range
}

func TestBuggyDispatch(t *testing.T) {
	type Foo struct {
		Id int `sql:"primary"`
	}
	dummy := NewDummyService()
	buggy := NewBuggyService()

	//A simple dispatch function
	buggyDispatch := func(record interface{}) (int, error) {
		return 0, RecordError("Dispatch has a bug")
	}

	service := NewDispatchService(buggyDispatch, []RecordService{dummy, buggy})
	err := service.Create(Foo{})
	if err.Error() != "Dispatch has a bug" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Update(Foo{})
	if err.Error() != "Dispatch has a bug" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Delete(Foo{})
	if err.Error() != "Dispatch has a bug" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Read(Foo{})
	if err.Error() != "Dispatch has a bug" {
		t.Errorf("The error message was wrong, got %v", err)
	}

	//A simple dispatch function
	outOfRange := func(record interface{}) (int, error) {
		return 2, nil
	}

	service = NewDispatchService(outOfRange, []RecordService{dummy, buggy})
	err = service.Create(Foo{})
	if err.Error() != "Dispatch index out of range" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Update(Foo{})
	if err.Error() != "Dispatch index out of range" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Delete(Foo{})
	if err.Error() != "Dispatch index out of range" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Read(Foo{})
	if err.Error() != "Dispatch index out of range" {
		t.Errorf("The error message was wrong, got %v", err)
	}
}
