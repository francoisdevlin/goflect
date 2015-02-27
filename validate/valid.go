package validate

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
)

func DefaultMatcher(record interface{}) (matcher.Matcher, error) {
	output := matcher.Any()
	fields := goflect.GetInfo(record)
	kinds := make(map[string]reflect.Kind)
	for _, field := range fields {
		kinds[field.Name] = field.Kind
	}
	p := matcher.ParseStruct{Fields: kinds}
	for _, field := range fields {
		fmt.Println(field)
		match, err := p.Parse("")
		//match, err := p.Parse(field.ValidExpr)
		if err != nil {
			return nil, err
		}
		output = matcher.And(output, match)
	}
	return output, nil
}
