package records

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"testing"
)

func TestDummyService(t *testing.T) {
	dummy := dummyService{}
	dummy.Insert(&dummy)
	dummy.ReadAll(&dummy)
	dummy.Update(&dummy)
	dummy.Delete(&dummy)
	if (dummy != dummyService{Inserts: 1, Updates: 1, Reads: 1, Deletes: 1}) {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}
}

/*
This shows how to use the a service to resetrict behavior
*/
func ExampleRecordService_basicRestrict() {
	match := matcher.StructMatcher{}
	match.AddField("Inserts", matcher.Gte(1))

	dummy := NewDummyService()
	service, _ := dummy.Restrict(match)
	//Inserts are blocked
	err := service.Insert(dummyService{})
	fmt.Println(err)
	//Output:
	//Could not insert record, does not match
}
