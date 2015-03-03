package records

import (
	"fmt"
	//"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
This method requires a primary key, and will return an error when not provided
*/
func ExampleRecordService_Delete_primaryKeyRequired() {
	type Foo struct {
		A int
	}
	service := NewDummyService()

	err := service.Delete(&Foo{})
	fmt.Println(err)

	//Output:
	//Bacon
}

/*
This method requires a primary key, and will return an error when not provided
*/
func ExampleRecordService_Update_primaryKeyRequired() {
	type Foo struct {
		A int
	}
	service := NewDummyService()

	err := service.Update(&Foo{})
	fmt.Println(err)

	//Output:
	//Bacon
}
