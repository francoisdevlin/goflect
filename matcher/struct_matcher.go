package matcher

import (
	"reflect"
)

type fieldYielder struct {
	matcher *structMatcher
	Name    string
}

func lookup(record interface{}, name string) (interface{}, error) {
	switch r := record.(type) {
	case map[string]interface{}:
		output, ok := r[name]
		if !ok {
			return nil, InvalidCompare(2)
		}
		return output, nil
	default:
		val := reflect.ValueOf(r)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if !val.IsValid() {
			return nil, InvalidCompare(2)
		}
		rVal := val.FieldByName(name)
		if !rVal.IsValid() {
			return nil, InvalidCompare(2)
		}
		return rVal.Interface(), nil
	}
}

func (y fieldYielder) Yield() (interface{}, error) {
	record := y.matcher.record
	return lookup(record, y.Name)
}

/*
This type is the main item to use with the matcher API
*/
type structMatcher struct {
	record interface{} //This is used to make the yielder work
	Fields map[string]Matcher
}

func (field *structMatcher) AddField(name string, matcher Matcher) {
	if field.Fields == nil {
		field.Fields = make(map[string]Matcher)
	}
	field.Fields[name] = matcher
}

func (field *structMatcher) Field(Name string) Yielder {
	return fieldYielder{Name: Name, matcher: field}
}

func (field *structMatcher) Match(record interface{}) (bool, error) {
	if field.Fields == nil {
		field.Fields = make(map[string]Matcher)
	}
	field.record = record
	valid := true
	for name, matcher := range field.Fields {
		attr, err := lookup(record, name)
		if err != nil {
			return false, InvalidCompare(1)
		}
		localMatch, err := matcher.Match(attr)
		if err != nil {
			return localMatch, err
		}
		valid = valid && localMatch
		if !valid {
			return false, nil
		}
	}
	return true, nil
}
