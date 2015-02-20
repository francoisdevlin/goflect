package goflect

import (
	"fmt"
	"sort"
	"strings"
)

type Printer interface {
	Print(m Matcher) (string, error)
}

type DefaultPrinter struct {
	Var string
}

/*
This prints a human readable representation of the matcher.  It is
*/
func (p DefaultPrinter) Print(m Matcher) (string, error) {
	switch r := m.(type) {
	case AndMatch:
		output := make([]string, 0)
		for _, matcher := range r.Matchers {
			result, err := p.Print(matcher)
			if err != nil {
				return "", err
			}
			result = "(" + result + ")"
			output = append(output, result)
		}
		return strings.Join(output, " AND "), nil
	case OrMatch:
		output := make([]string, 0)
		for _, matcher := range r.Matchers {
			result, err := p.Print(matcher)
			if err != nil {
				return "", err
			}
			result = "(" + result + ")"
			output = append(output, result)
		}
		return strings.Join(output, " OR "), nil
	case StructMatcher:
		output := make([]string, 0)
		keys := make([]string, 0)
		for name, _ := range r.Fields {
			keys = append(keys, name)
		}
		sort.Strings(keys)
		for _, name := range keys {
			matcher, _ := r.Fields[name]
			printer := DefaultPrinter{Var: name}
			result, err := printer.Print(matcher)
			if err != nil {
				return "", err
			}
			result = "(" + result + ")\n"
			output = append(output, result)
		}
		return strings.Join(output, "AND "), nil
	case FieldMatcher:
		output := ""
		if p.Var == "" {
			output += "_"
		} else {
			output += p.Var
		}
		output += " "
		switch r.Op {
		case EQ:
			output += "="
		case NEQ:
			output += "!="
		case LT:
			output += "<"
		case LTE:
			output += "<="
		case GT:
			output += ">"
		case GTE:
			output += ">="
		case IN:
			output += "IN"
		case NOT_IN:
			output += "NOT IN"
		case MATCH:
			output += "MATCH"
		case NOT_MATCH:
			output += "NOT MATCH"
		}
		switch val := r.Value.(type) {
		case []string:
			output += " ["
			entries := make([]string, 0)
			for _, v := range val {
				entries = append(entries, "\""+v+"\"")
			}
			output += strings.Join(entries, " ")
			output += "]"
			return output, nil
		case string:
			return output + " \"" + val + "\"", nil
		default:
			return output + " " + fmt.Sprint(r.Value), nil
		}
	case InvertMatch:
		recurse, err := p.Print(r.M)
		if err != nil {
			return "", err
		}
		return "NOT (" + recurse + ")", nil
	case NoneMatch:
		return "false", nil
	case AnyMatch:
		return "true", nil
	}
	return "", nil
}
