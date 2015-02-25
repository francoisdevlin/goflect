package matcher

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Printer interface {
	Print(m Matcher) (string, error)
}

type Literal string
type DefaultPrinter struct {
	Var string
}

type SqlitePrinter struct {
	Var string
}

func printAnd(p Printer, r AndMatch) (string, error) {
	output := make([]string, 0)
	for _, matcher := range r.Matchers {
		result, err := p.Print(matcher)
		if err != nil {
			return "", err
		}
		needsParens := false
		switch temp := matcher.(type) {
		case OrMatch:
			needsParens = true
		default:
			temp = temp
		}
		if needsParens {
			result = "(" + result + ")"
		}
		output = append(output, result)
	}
	return strings.Join(output, " AND "), nil
}

func printOr(p Printer, r OrMatch) (string, error) {
	output := make([]string, 0)
	for _, matcher := range r.Matchers {
		result, err := p.Print(matcher)
		if err != nil {
			return "", err
		}
		needsParens := false
		switch temp := matcher.(type) {
		case AndMatch:
			needsParens = true
		case StructMatcher:
			needsParens = true
		default:
			temp = temp
		}
		if needsParens {
			result = "(" + result + ")"
		}
		output = append(output, result)
	}
	return strings.Join(output, " OR "), nil
}

func printInvert(p Printer, r InvertMatch) (string, error) {
	recurse, err := p.Print(r.M)
	if err != nil {
		return "", err
	}
	return "NOT (" + recurse + ")", nil
}

func printStruct(factory func(name string) Printer, r StructMatcher) (string, error) {
	output := make([]string, 0)
	keys := make([]string, 0)
	for name, _ := range r.Fields {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	for _, name := range keys {
		matcher, _ := r.Fields[name]
		printer := factory(name)
		result, err := printer.Print(matcher)
		if err != nil {
			return "", err
		}
		needsParens := false
		switch temp := matcher.(type) {
		case OrMatch:
			needsParens = true
		default:
			temp = temp
		}
		if needsParens {
			result = "(" + result + ")"
		}
		output = append(output, result)
	}
	return strings.Join(output, " AND "), nil
}

func q(name string) string {
	return "\"" + name + "\""
}

func stringToken(token FieldOps) string {
	output := ""
	switch token {
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
	return output
}

/*
This prints a human readable representation of the matcher.  It is
*/
func (p DefaultPrinter) Print(m Matcher) (string, error) {
	switch r := m.(type) {
	case AndMatch:
		return printAnd(p, r)
	case OrMatch:
		return printOr(p, r)
	case StructMatcher:
		return printStruct(func(name string) Printer { return DefaultPrinter{Var: name} }, r)
	case FieldMatcher:
		output := ""
		if p.Var == "" {
			output += "_"
		} else {
			output += p.Var
		}
		output += " " + stringToken(r.Op)
		switch val := r.Value.(type) {
		case []string:
			output += " ["
			entries := make([]string, 0)
			for _, v := range val {
				entries = append(entries, q(v))
			}
			output += strings.Join(entries, " ")
			output += "]"
			return output, nil
		case string:
			return output + " " + q(val), nil
		case Literal:
			return output + " " + string(val), nil
		default:
			return output + " " + fmt.Sprint(r.Value), nil
		}
	case InvertMatch:
		return printInvert(p, r)
	case NoneMatch:
		return "false", nil
	case AnyMatch:
		return "true", nil
	}
	return "", nil
}

/*
This prints a human readable representation of the matcher.  It is
*/
func (p SqlitePrinter) Print(m Matcher) (string, error) {
	switch r := m.(type) {
	case AndMatch:
		return printAnd(p, r)
	case OrMatch:
		return printOr(p, r)
	case StructMatcher:
		return printStruct(func(name string) Printer { return SqlitePrinter{Var: name} }, r)
	case FieldMatcher:
		output := ""
		if p.Var == "" {
			output += "_"
		} else {
			output += p.Var
		}
		output += " " + stringToken(r.Op)
		makeInish := func(entries []string) string {
			output += " ("
			output += strings.Join(entries, ", ")
			output += ")"
			return output
		}
		entries := make([]string, 0)
		switch val := r.Value.(type) {
		case []int:
			for _, v := range val {
				entries = append(entries, strconv.FormatInt(int64(v), 10))
			}
			return makeInish(entries), nil
		case []int64:
			for _, v := range val {
				entries = append(entries, strconv.FormatInt(v, 10))
			}
			return makeInish(entries), nil
		case []int32:
			for _, v := range val {
				entries = append(entries, strconv.FormatInt(int64(v), 10))
			}
			return makeInish(entries), nil
		case []int16:
			for _, v := range val {
				entries = append(entries, strconv.FormatInt(int64(v), 10))
			}
			return makeInish(entries), nil
		case []int8:
			for _, v := range val {
				entries = append(entries, strconv.FormatInt(int64(v), 10))
			}
			return makeInish(entries), nil
		case []uint:
			for _, v := range val {
				entries = append(entries, strconv.FormatUint(uint64(v), 10))
			}
			return makeInish(entries), nil
		case []uint64:
			for _, v := range val {
				entries = append(entries, strconv.FormatUint(v, 10))
			}
			return makeInish(entries), nil
		case []uint32:
			for _, v := range val {
				entries = append(entries, strconv.FormatUint(uint64(v), 10))
			}
			return makeInish(entries), nil
		case []uint16:
			for _, v := range val {
				entries = append(entries, strconv.FormatUint(uint64(v), 10))
			}
			return makeInish(entries), nil
		case []uint8:
			for _, v := range val {
				entries = append(entries, strconv.FormatUint(uint64(v), 10))
			}
			return makeInish(entries), nil
		case []float64:
			for _, v := range val {
				entries = append(entries, strconv.FormatFloat(v, 'f', 4, 64))
			}
			return makeInish(entries), nil
		case []float32:
			for _, v := range val {
				entries = append(entries, strconv.FormatFloat(float64(v), 'f', 4, 32))
			}
			return makeInish(entries), nil
		case []string:
			for _, v := range val {
				entries = append(entries, q(v))
			}
			return makeInish(entries), nil
		case string:
			return output + " " + q(val), nil
		case Literal:
			return output + " " + string(val), nil
		default:
			return output + " " + fmt.Sprint(r.Value), nil
		}
	case InvertMatch:
		return printInvert(p, r)
	case NoneMatch:
		return "false", nil
	case AnyMatch:
		return "true", nil
	}
	return "", nil
}
