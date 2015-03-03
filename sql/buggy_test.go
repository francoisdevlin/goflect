package records

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
This demonstrates that the buggy service can be used to generate errors
*/
func ExampleNewBuggyService_1() {
	type Foo struct {
		A int `sql:"primary"`
	}
	service := NewBuggyService()
	//Create Errors
	err := service.Create(&Foo{})
	fmt.Println(err)
	//Update Errors
	err = service.Update(&Foo{})
	fmt.Println(err)
	err = service.UpdateAll(&Foo{})
	fmt.Println(err)
	err = service.UpdateAllWhere(&Foo{}, matcher.Any())
	fmt.Println(err)
	//Read Errors
	_, err = service.ReadAll(&Foo{})
	fmt.Println(err)
	_, err = service.ReadAllWhere(&Foo{}, matcher.Any())
	fmt.Println(err)
	err = service.Get(0, &Foo{})
	fmt.Println(err)
	//Delete Errors
	err = service.Delete(&Foo{})
	fmt.Println(err)
	err = service.DeleteAll(&Foo{})
	fmt.Println(err)
	err = service.DeleteAllWhere(&Foo{}, matcher.Any())
	fmt.Println(err)
	err = service.DeleteById(0, &Foo{})
	fmt.Println(err)

	//Output:
	//Intentional Create Error
	//Intentional Update Error
	//Intentional Update Error
	//Intentional Update Error
	//Intentional Read Error
	//Intentional Read Error
	//Intentional Read Error
	//Intentional Delete Error
	//Intentional Delete Error
	//Intentional Delete Error
	//Intentional Delete Error
}
