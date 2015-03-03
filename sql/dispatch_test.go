package records

import (
	"fmt"
	//"testing"
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
