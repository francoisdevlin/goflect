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
	for _, field := range fields {
		fieldVal := val.FieldByName(field.Name)
		if field.IsPrimary {
			match.AddField(field.Name, matcher.Eq(fieldVal.Interface()))
		}
	}
	return match, nil
}

func (service RecordService) Delete(record interface{}) error {
	match, _ := primaryMatcher(record)
	service.delegate.deleteAll(record, match)
	return nil
}

func (service RecordService) Update(record interface{}) error {
	match, _ := primaryMatcher(record)
	service.delegate.updateAll(record, match)
	return nil
}

func (service RecordService) Insert(record interface{}) error {
	sliceType := reflect.SliceOf(reflect.TypeOf(record))
	slice := reflect.MakeSlice(sliceType, 0, 1)
	slice = reflect.Append(slice, reflect.ValueOf(record))
	return service.delegate.insertAll(slice.Interface())
}

func (service RecordService) Get(id int64, record interface{}) {
	fields := goflect.GetInfo(record)
	match := matcher.NewStructMatcher()
	for _, field := range fields {
		if field.IsPrimary {
			match.AddField(field.Name, matcher.Eq(id))
		}
	}
	next, _ := service.ReadWhere(record, match)
	for next(record) {
	} //The last call closes the result set, important!
}

func (service RecordService) DeleteById(id int64, record interface{}) {
	fields := goflect.GetInfo(record)
	match := matcher.NewStructMatcher()
	for _, field := range fields {
		if field.IsPrimary {
			match.AddField(field.Name, matcher.Eq(id))
		}
	}
	_ = service.delegate.deleteAll(record, match)
}

/*
This function can be used to return a set of records that match a set of criteria.  It accepts a matcher that describes a record set.
*/
func (service RecordService) ReadWhere(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
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