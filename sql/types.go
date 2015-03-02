/*
This is a package for abstracting connection to record sources.  It is designed to do most of the heavy lifting in an application for the user.
*/
package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

type RecordDefiner interface {
	Create(record interface{}) error
}

type SqlRecordDefiner interface {
	RecordDefiner
	CreateStatement(record interface{}) string
}

type RecordError string

func (e RecordError) Error() string {
	return string(e)
}

type Nominal struct {
	Id   int64
	Name string
}

type privateRecordService interface {
	insertAll(rows interface{}) error
	readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error)
	updateAll(record interface{}, match matcher.Matcher) error
	deleteAll(record interface{}, match matcher.Matcher) error
}

/*
This is yet another take on the active record pattern, centered around the matcher.  It also includes also sort of distributed computing hotness by leveraging privateRecordServices
*/
type RecordService struct {
	delegate privateRecordService
}

//ReadAllNominalWhere(record interface{}, conditions map[string]interface{}) func(record *Nominal) bool
//GetByNominal(name string, record interface{})
//DeleteById(id int64, record interface{})
//GetNominalByNominal(name string) (output Nominal)
//Get(id int64, record interface{})
//GetNominal(id int64) (output Nominal)
//ReadAllNominal(record interface{}) (func(record *Nominal) bool, error)
//Limit(int) (RecordService error)
//Offset(int) (RecordService error)
//OrderBy(...string) (RecordService error)
