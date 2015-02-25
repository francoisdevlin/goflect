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
		dum1, _ := d1.(*dummyService)
		dum2, _ := d2.(*dummyService)
		fmt.Println("Inserts - dummy 1, dummy 2:", dum1.Inserts, dum2.Inserts)
	}

	//Create our dummy services
	dummy1 := NewDummyService()
	dummy2 := NewDummyService()

	//A simple dispatch function
	disp := func(record interface{}) (int, error) {
		dummy, _ := record.(dummyService)
		return dummy.Inserts % 2, nil
	}

	//Create the dispatch service
	service := NewDispatchService(disp, []RecordService{dummy1, dummy2})

	//We start with nothing
	printStats(dummy1, dummy2)

	//The insert only goes to one dummy
	service.Insert(dummyService{Inserts: 1})
	printStats(dummy1, dummy2)

	//The insert goes to the other dummy
	service.Insert(dummyService{Inserts: 2})
	printStats(dummy1, dummy2)
	//Output:
	//Inserts - dummy 1, dummy 2: 0 0
	//Inserts - dummy 1, dummy 2: 0 1
	//Inserts - dummy 1, dummy 2: 1 1
}
