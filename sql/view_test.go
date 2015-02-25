package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"testing"
)

func TestView(t *testing.T) {
	match := matcher.StructMatcher{}
	match.AddField("Inserts", matcher.Gte(1))

	dummy := DummyService{}
	service, _ := dummy.Restrict(match)
	//Inserts are blocked
	service.Insert(DummyService{})
	if (dummy != DummyService{}) {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}

	//Updates are blocked
	service.Update(DummyService{})
	if (dummy != DummyService{}) {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}

	//Deletes are blocked
	service.Delete(DummyService{})
	if (dummy != DummyService{}) {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}

	//The service methods work properly when the record meets the conditions
	oneInsert := DummyService{Inserts: 1}
	oneAll := DummyService{Inserts: 1, Updates: 1, Reads: 1, Deletes: 1}
	service.Insert(oneInsert)
	service.ReadAll(oneInsert)
	service.Update(oneInsert)
	service.Delete(oneInsert)
	if dummy != oneAll {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}

	match = matcher.StructMatcher{}
	match.AddField("Updates", matcher.Gte(1))
	service, _ = service.Restrict(match)

	//Inserts are blocked
	service.Insert(oneInsert)
	if dummy != oneAll {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}

	//Updates are blocked
	service.Update(oneInsert)
	if dummy != oneAll {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}

	//Deletes are blocked
	service.Delete(oneInsert)
	if dummy != oneAll {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}

	service.Insert(oneAll)
	service.ReadAll(oneAll)
	service.Update(oneAll)
	service.Delete(oneAll)
	if (dummy != DummyService{Inserts: 2, Updates: 2, Reads: 2, Deletes: 2}) {
		t.Errorf("The dummy is not the expected value, it is: %v", dummy)
	}

}
