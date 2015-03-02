package matcher

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func NewSqlitePrinter() Printer {
	return sqlitePrinter{}
}

func NewDefaultPrinter() Printer {
	return defaultPrinter{}
}

type defaultPrinter struct {
	v string
}

type sqlitePrinter struct {
	v string
}

func printAnd(p Printer, r andMatch) (string, error) {
	output := make([]string, 0)
	for _, matcher := range r.Matchers {
		result, err := p.Print(matcher)
		if err != nil {
			return "", err
		}
		needsParens := false
		switch temp := matcher.(type) {
		case orMatch:
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

func printOr(p Printer, r orMatch) (string, error) {
	output := make([]string, 0)
	for _, matcher := range r.Matchers {
		result, err := p.Print(matcher)
		if err != nil {
			return "", err
		}
		needsParens := false
		switch temp := matcher.(type) {
		case andMatch:
			needsParens = true
		case *structMatcher:
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

func printInvert(p Printer, r invertMatch) (string, error) {
	recurse, err := p.Print(r.M)
	if err != nil {
		return "", err
	}
	return "NOT (" + recurse + ")", nil
}

func printStruct(factory func(name string) Printer, r *structMatcher) (string, error) {
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
		case orMatch:
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

//func stringToken(token fieldOps) string {
func (token fieldOps) String() string {
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
func (p defaultPrinter) Print(m Matcher) (string, error) {
	switch r := m.(type) {
	case andMatch:
		return printAnd(p, r)
	case orMatch:
		return printOr(p, r)
	case *structMatcher:
		return printStruct(func(name string) Printer { return defaultPrinter{v: name} }, r)
	case fieldMatcher:
		output := ""
		if p.v == "" {
			output += "_"
		} else {
			output += p.v
		}
		output += " " + r.Op.String()
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
		case fieldYielder:
			return output + " " + val.Name, nil
		default:
			return output + " " + fmt.Sprint(r.Value), nil
		}
	case invertMatch:
		return printInvert(p, r)
	case noneMatch:
		return "false", nil
	case anyMatch:
		return "true", nil
	case errorMatch:
		return "error", nil
	}
	return "", nil
}

/*
This prints a human readable representation of the matcher.  It is specifically tweaks to provide a vlaid where clause for SQLite
*/
func (p sqlitePrinter) Print(m Matcher) (string, error) {
	switch r := m.(type) {
	case andMatch:
		return printAnd(p, r)
	case orMatch:
		return printOr(p, r)
	case *structMatcher:
		return printStruct(func(name string) Printer { return sqlitePrinter{v: name} }, r)
	case fieldMatcher:
		output := ""
		if p.v == "" {
			output += "_"
		} else {
			output += p.v
		}
		output += " " + r.Op.String()
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
		case fieldYielder:
			return output + " " + val.Name, nil
		default:
			return output + " " + fmt.Sprint(r.Value), nil
		}
	case invertMatch:
		return printInvert(p, r)
	case noneMatch:
		return "0", nil
	case anyMatch:
		return "1", nil
	case errorMatch:
		return "", InvalidCompare(0)
	}
	return "", nil
}
