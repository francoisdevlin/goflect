package matcher

import (
	//"fmt"
	"reflect"
	//"regexp"
)

type FieldYielder struct {
	matcher *StructMatcher
	Name    string
}

func (y FieldYielder) Yield() (interface{}, error) {
	record := y.matcher.record
	switch r := record.(type) {
	case map[string]interface{}:
		output, ok := r[y.Name]
		if !ok {
			return nil, InvalidCompare(2)
		}
		return output, nil
	default:
		val := reflect.ValueOf(r)
		if !val.IsValid() {
			return nil, InvalidCompare(2)
		}
		rVal := val.FieldByName(y.Name)
		if !rVal.IsValid() {
			return nil, InvalidCompare(2)
		}
		return rVal.Interface(), nil
	}
	return nil, InvalidCompare(2)
}

/*
This type is the main item to use with the matcher API
*/
type StructMatcher struct {
	record interface{} //This is used to make the yielder work
	Fields map[string]Matcher
}

func (field *StructMatcher) AddField(name string, matcher Matcher) {
	if field.Fields == nil {
		field.Fields = make(map[string]Matcher)
	}
	field.Fields[name] = matcher
}

func (field *StructMatcher) Field(Name string) Yielder {
	return FieldYielder{Name: Name, matcher: field}
}

func (field *StructMatcher) Match(record interface{}) (bool, error) {
	if field.Fields == nil {
		field.Fields = make(map[string]Matcher)
	}
	field.record = record
	switch r := record.(type) {
	case map[string]interface{}:
		valid := true
		for name, matcher := range field.Fields {
			attr, present := r[name]
			if !present {
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
	default:
		val := reflect.ValueOf(record)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		valid := true

		for name, matcher := range field.Fields {
			attr := val.FieldByName(name)
			if !attr.IsValid() {
				return false, InvalidCompare(1)
			}
			localMatch, err := matcher.Match(attr.Interface())
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

	return false, InvalidCompare(1)
}
