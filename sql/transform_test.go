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
		dum1, _ := d1.(*dummyService)
		fmt.Println("Inserts - dummy:", dum1.Inserts)
	}

	//A simple match
	match := matcher.NewStructMatcher()
	match.AddField("Inserts", matcher.Gte(1))

	dummy := NewDummyService()
	service, _ := dummy.Restrict(match)
	//Inserts are blocked
	err := service.Insert(dummyService{})
	fmt.Println(err)
	printStats(dummy)

	//Our transform function
	insertTransform := func(record interface{}) (interface{}, error) {
		dummy, _ := record.(dummyService)
		dummy.Inserts++
		return dummy, nil
	}
	//But, if we add a transform to the service...
	service = NewTransformService(insertTransform, service)
	service.Insert(dummyService{})
	printStats(dummy)

	//Output:
	//Could not insert record, does not match
	//Inserts - dummy: 0
	//Inserts - dummy: 1
}
