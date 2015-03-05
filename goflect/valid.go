package goflect

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
)

/*
The default matcher is used to determine the default matcher for a given type.  It uses the information in the valid tag in order to do this.
*/
func DefaultMatcher(record interface{}) (matcher.Matcher, error) {
	output := matcher.Any()
	fields := GetInfo(record)
	p, err := DefaultParser(record)
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		match, err := p.Parse(field.ValidExpr)
		if err != nil {
			return nil, err
		}
		output = matcher.And(output, match)
	}
	return output, nil
}

/*
This uses the provided record to create a context for the parser.  It the record is a struct, reflection is used.  If the record is a primitive, its type is used and a primitive parser is return instead
*/
func DefaultParser(record interface{}) (matcher.Parser, error) {
	typ := reflect.TypeOf(record)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	kind := typ.Kind()

	if kind == reflect.Interface || kind == reflect.Struct {

		fields := GetInfo(record)
		kinds := make(map[string]reflect.Kind)
		for _, field := range fields {
			kinds[field.Name] = field.Kind
		}
		return matcher.NewParser(kinds)
	} else {
		return matcher.NewParser(record)
	}
}

/*
This uses default parser to parse the input string, and provide a matcher
*/
func Parse(record interface{}, input string) (matcher.Matcher, error) {
	p, err := DefaultParser(record)
	if err != nil {
		return nil, err
	}
	return p.Parse(input)
}
