/*
The plan is to use matcher objects to define projections in the product.  This can be used in multiple ways

    A composable SQL querying API
    A filtering mechansim
    A way of validating data is valid for a struct
    A way of providing a filtering mechanism for user interaction
    A consistent rules API mechanism
*/
package goflect

import (
	//"fmt"
	"reflect"
	//"regexp"
)

/*
This type is the main item to use with the matcher API
*/
type StructMatcher struct {
	Fields map[string]Matcher
}

func (field *StructMatcher) AddField(name string, matcher Matcher) {
	if field.Fields == nil {
		field.Fields = make(map[string]Matcher)
	}
	field.Fields[name] = matcher
}

func (field StructMatcher) Match(record interface{}) (bool, error) {
	if field.Fields == nil {
		field.Fields = make(map[string]Matcher)
	}
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
