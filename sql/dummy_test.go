package records

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"testing"
)

func TestDummyService(t *testing.T) {
	service := NewDummyService()
	dummy, _ := service.delegate.(*dummyService)
	service.Create(dummy)
	service.ReadAll(dummy)
	service.UpdateAll(dummy)
	service.DeleteAll(dummy)
	if (*dummy != dummyService{Creates: 1, Updates: 1, Reads: 1, Deletes: 1}) {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}
}

/*
This shows how to use the a service to resetrict behavior
*/
func ExampleRecordService_basicRestrict() {
	match := matcher.NewStructMatcher()
	match.AddField("Creates", matcher.Gte(1))

	dummy := NewDummyService()
	service, _ := NewViewService(match, dummy)
	//Creates are blocked
	err := service.Create(dummyService{})
	fmt.Println(err)
	//Output:
	//Could not create record, does not match
}
