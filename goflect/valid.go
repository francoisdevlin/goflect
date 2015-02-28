package goflect

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
)

func DefaultMatcher(record interface{}) (matcher.Matcher, error) {
	output := matcher.Any()
	fields := GetInfo(record)
	kinds := make(map[string]reflect.Kind)
	for _, field := range fields {
		kinds[field.Name] = field.Kind
	}
	p := matcher.ParseStruct{Fields: kinds}
	for _, field := range fields {
		match, err := p.Parse(field.ValidExpr)
		if err != nil {
			return nil, err
		}
		output = matcher.And(output, match)
	}
	return output, nil
}
