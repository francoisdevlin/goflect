package records

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"testing"
)

/*
This shows how to use the a service to resetrict behavior
*/
func ExampleRecordService_basicTransform() {
	//A basic printer
	printStats := func(d1 RecordService) {
		dum1, _ := d1.delegate.(*dummyService)
		fmt.Println("Creates - dummy:", dum1.Creates)
	}

	//A simple match
	match := matcher.NewStructMatcher()
	match.AddField("Creates", matcher.Gte(1))

	dummy := NewDummyService()
	service, _ := NewViewService(match, dummy)
	//Creates are blocked
	err := service.Create(dummyService{})
	fmt.Println(err)
	printStats(dummy)

	//Our transform function
	createTransform := func(record interface{}) (interface{}, error) {
		dummy, _ := record.(dummyService)
		dummy.Creates++
		return dummy, nil
	}
	////But, if we add a transform to the service...
	service = NewTransformService(createTransform, service)
	service.Create(dummyService{})
	printStats(dummy)

	//Output:
	//Could not create record, does not match
	//Creates - dummy: 0
	//Creates - dummy: 1
}

func TestBuggyTransform(t *testing.T) {
	type Foo struct {
		Id int `sql:"primary"`
	}
	dummy := NewDummyService()

	//A simple dispatch function
	buggyTransform := func(record interface{}) (interface{}, error) {
		return nil, RecordError("Transform has a bug")
	}

	service := NewTransformService(buggyTransform, dummy)
	err := service.Create(Foo{})
	if err.Error() != "Transform has a bug" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Update(Foo{})
	if err.Error() != "Transform has a bug" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Delete(Foo{})
	if err.Error() != "Transform has a bug" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Read(Foo{})
	if err.Error() != "Transform has a bug" {
		t.Errorf("The error message was wrong, got %v", err)
	}

	//A simple dispatch function
	nilTransform := func(record interface{}) (interface{}, error) {
		return nil, nil
	}

	service = NewTransformService(nilTransform, dummy)
	err = service.Create(Foo{})
	if err.Error() != "Tranform returned nil" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Update(Foo{})
	if err.Error() != "Tranform returned nil" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Delete(Foo{})
	if err.Error() != "Tranform returned nil" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Read(Foo{})
	if err.Error() != "Tranform returned nil" {
		t.Errorf("The error message was wrong, got %v", err)
	}

	//A simple dispatch function
	identityTransform := func(record interface{}) (interface{}, error) {
		return record, nil
	}

	service = NewTransformService(identityTransform, NewBuggyService())
	err = service.Create(Foo{})
	if err.Error() != "Intentional Create Error" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Update(Foo{})
	if err.Error() != "Intentional Update Error" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Delete(Foo{})
	if err.Error() != "Intentional Delete Error" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Read(Foo{})
	if err.Error() != "Intentional Read Error" {
		t.Errorf("The error message was wrong, got %v", err)
	}
}
