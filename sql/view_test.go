package records

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"testing"
)

func TestView(t *testing.T) {
	match := matcher.NewStructMatcher()
	match.AddField("Creates", matcher.Gte(1))

	service := NewDummyService()
	dummy, _ := service.delegate.(*dummyService)
	service, _ = NewViewService(match, service)
	//Creates are blocked
	service.Create(dummyService{})
	if (*dummy != dummyService{}) {
		t.Errorf("1. The dummy is not the expected value, it is: %v", dummy)
	}

	//Updates are blocked
	service.UpdateAll(dummyService{})
	if (*dummy != dummyService{}) {
		t.Errorf("2. The dummy is not the expected value, it is: %v", dummy)
	}

	//Deletes are blocked
	service.DeleteAll(dummyService{})
	if (*dummy != dummyService{}) {
		t.Errorf("3. The dummy is not the expected value, it is: %v", dummy)
	}

	//The service methods work properly when the record meets the conditions
	oneCreate := dummyService{Creates: 1}
	oneAll := dummyService{Creates: 1, Updates: 1, Reads: 1, Deletes: 1}
	service.Create(oneCreate)
	service.ReadAll(oneCreate)
	service.UpdateAll(oneCreate)
	service.DeleteAll(oneCreate)
	if *dummy != oneAll {
		t.Errorf("4. The dummy is not the expected value, it is: %v", dummy)
	}

	match = matcher.NewStructMatcher()
	match.AddField("Updates", matcher.Gte(1))
	service, _ = NewViewService(match, service)

	//Creates are blocked
	service.Create(oneCreate)
	if *dummy != oneAll {
		t.Errorf("5. The dummy is not the expected value, it is: %v", dummy)
	}

	//Updates are blocked
	service.UpdateAll(oneCreate)
	if *dummy != oneAll {
		t.Errorf("6. The dummy is not the expected value, it is: %v", dummy)
	}

	//Deletes are blocked
	service.DeleteAll(oneCreate)
	if *dummy != oneAll {
		t.Errorf("7. The dummy is not the expected value, it is: %v", dummy)
	}

	service.Create(oneAll)
	service.ReadAll(oneAll)
	service.UpdateAll(oneAll)
	service.DeleteAll(oneAll)
	if (*dummy != dummyService{Creates: 2, Updates: 2, Reads: 2, Deletes: 2}) {
		t.Errorf("8. The dummy is not the expected value, it is: %v", dummy)
	}

}

func TestBuggyMatcher(t *testing.T) {
	type Foo struct {
		Id int `sql:"primary"`
	}
	dummy := NewDummyService()

	service, _ := NewViewService(matcher.Buggy(), dummy)
	err := service.Create(Foo{})
	if err.Error() != "Invalid comparison operation" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Update(Foo{})
	if err.Error() != "Invalid comparison operation" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Delete(Foo{})
	if err.Error() != "Invalid comparison operation" {
		t.Errorf("The error message was wrong, got %v", err)
	}
	err = service.Read(Foo{})
	if err.Error() != "Invalid comparison operation" {
		t.Errorf("The error message was wrong, got %v", err)
	}

	service, _ = NewViewService(matcher.Any(), NewBuggyService())
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

func ExampleNewViewService_deleteRecords() {
	//Start with our service boilerplate, record insertion
	service := tableCreationBoilerplate()
	service.Create(Device{Name: "Device 1"})
	service.Create(Device{Name: "Device 2"})
	service.Create(Device{Name: "Device 3"})
	service.Create(Device{Name: "Device 4"})

	//Create a matcher that looks for Ids greater than 2
	constraint := matcher.NewStructMatcher()
	constraint.AddField("Id", matcher.Gt(int64(2)))

	//Create the view service
	viewService, _ := NewViewService(constraint, service)
	device := Device{}
	fmt.Println("Printing devices in view")
	printAll(viewService, &device)

	//Delete all of the devices in the view
	device = Device{Id: 3, Name: "Bacon"}
	viewService.DeleteAll(&device)
	fmt.Println("Printing devices in view")
	printAll(viewService, &device)

	//Print the remaining devices
	fmt.Println("Printing all remaining devices")
	printAll(service, &device)
	//Output:
	//Printing devices in view
	//Devices
	//{3 Device 3}
	//{4 Device 4}
	//Printing devices in view
	//Devices
	//Printing all remaining devices
	//Devices
	//{1 Device 1}
	//{2 Device 2}
}
