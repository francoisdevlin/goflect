package records

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	//"testing"
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
