package records

import (
	"database/sql"
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
)

func primaryMatcher(record interface{}) (matcher.Matcher, error) {
	_, val := typeAndVal(record)

	fields := goflect.GetInfo(record)
	match := matcher.NewStructMatcher()
	primaryFound := false
	for _, field := range fields {
		fieldVal := val.FieldByName(field.Name)
		if field.IsPrimary {
			match.AddField(field.Name, matcher.Eq(fieldVal.Interface()))
			primaryFound = true
		}
	}
	if primaryFound {
		return match, nil
	}
	return nil, RecordError("No primary key found")
}

/*
This method will delete the record specified by its primary key.  It will return an error if there is no primary key specified, or something goes wrong at a lower layer
*/
func (service RecordService) Delete(record interface{}) error {
	match, err := primaryMatcher(record)
	if err != nil {
		return err
	}
	return service.delegate.deleteAll(record, match)
}

/*
This method will delete all of the records associated with the record service
*/
func (service RecordService) DeleteAll(record interface{}) error {
	return service.delegate.deleteAll(record, matcher.Any())
}

/*
This method will delete all of the records associated with the record service, as restricted by the matcher
*/
func (service RecordService) DeleteAllWhere(record interface{}, match matcher.Matcher) error {
	return service.delegate.deleteAll(record, match)
}

/*
This method will update the record specified by its primary key.  It will return an error if there is no primary key specified, or something goes wrong at a lower layer
*/
func (service RecordService) Update(record interface{}) error {
	match, err := primaryMatcher(record)
	if err != nil {
		return err
	}
	return service.delegate.updateAll(record, match)
}

/*
This method will update all of the records associated with the record service to use the values in the record.  Really should only be combined with a dictionary, as to limit the fields that are updated
*/
func (service RecordService) UpdateAll(record interface{}) error {
	return service.delegate.updateAll(record, matcher.Any())
}

/*
This method will update all of the records associated with the record service to use the values in the record, as restricted by teh matcher.  Really should only be combined with a dictionary, as to limit the fields that are updated
*/
func (service RecordService) UpdateAllWhere(record interface{}, match matcher.Matcher) error {
	return service.delegate.updateAll(record, match)
}

/*
This creates a record into the service
*/
func (service RecordService) Create(record interface{}) error {
	sliceType := reflect.SliceOf(reflect.TypeOf(record))
	slice := reflect.MakeSlice(sliceType, 0, 1)
	slice = reflect.Append(slice, reflect.ValueOf(record))
	return service.delegate.createAll(slice.Interface())
}

func (service RecordService) Get(id int64, record interface{}) error {
	return service.ReadById(id, record)
}

func (service RecordService) ReadById(id int64, record interface{}) error {
	fields := goflect.GetInfo(record)
	match := matcher.NewStructMatcher()
	for _, field := range fields {
		if field.IsPrimary {
			match.AddField(field.Name, matcher.Eq(id))
		}
	}
	next, err := service.ReadAllWhere(record, match)
	if err != nil {
		return err
	}
	for next(record) {
	} //The last call closes the result set, important!
	return nil
}

func (service RecordService) DeleteById(id int64, record interface{}) error {
	fields := goflect.GetInfo(record)
	match := matcher.NewStructMatcher()
	for _, field := range fields {
		if field.IsPrimary {
			match.AddField(field.Name, matcher.Eq(id))
		}
	}
	return service.delegate.deleteAll(record, match)
}

/*
This function can be used to return a set of records that match a set of criteria.  It accepts a matcher that describes a record set.
*/
func (service RecordService) ReadAllWhere(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	return service.delegate.readAll(record, match)
}

/*
This returns all of the records that the service has access to
*/
func (service RecordService) ReadAll(record interface{}) (func(record interface{}) bool, error) {
	return service.delegate.readAll(record, matcher.Any())
}

/***
These are the constructors
***/

/*
This creates a new dispatch service that will route the request to the appropriate service underneath
*/
func NewDispatchService(disp func(interface{}) (int, error), delegs []RecordService) RecordService {
	services := make([]privateRecordService, len(delegs))
	for i, service := range delegs {
		services[i] = service.delegate
	}
	return RecordService{delegate: dispatch{dispatcher: disp, delegates: services}}
}

/*
This creates a new transform service that will route the request to the appropriate service underneath
*/
func NewTransformService(trans func(interface{}) (interface{}, error), deleg RecordService) RecordService {
	return RecordService{delegate: transform{transformer: trans, delegate: deleg.delegate}}
}

/*
This creates a new dummy to use for testing purposes
*/
func NewDummyService() RecordService {
	return RecordService{delegate: new(dummyService)}
}

/*
This creates a new view based on the underlying service
*/
func NewViewService(match matcher.Matcher, delegate RecordService) (RecordService, error) {
	return RecordService{delegate: view{match: match, delegate: delegate.delegate}}, nil
}

/*
This function takes a connection to a sqlite service, and returns a Record Service.  Should only be used with application setup code
*/
func NewSqliteService(conn *sql.DB) RecordService {
	return RecordService{delegate: sqliteRecordService{Conn: conn}}
}

/*
This returns a new buggy serivce, that always returns an error.  Useful for testing
*/
func NewBuggyService() RecordService {
	return RecordService{delegate: buggyService{}}
}
