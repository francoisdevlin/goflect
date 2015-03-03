package records

import (
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
