package lint

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"go/token"
	"regexp"
	"strconv"
	"strings"
)

func FormatStructTag(pos token.Position, input string) (string, []error) {
	backquotes, _ := regexp.Compile("(^`|`$)")
	input = backquotes.ReplaceAllString(input, "")
	tagKeys, errors := ParseStructTag(input)
	if len(errors) > 0 {
		fmt.Println(errors)
		return "", errors
	}
	cols := pos.Column - 1
	entries := make([]string, 0)
	touchedTags := make(map[string]int)
	fieldFormatter := map[string]func(string) string{
		goflect.TAG_SQL: FlagOrderFactory(goflect.SQL_FIELDS),
		goflect.TAG_UI:  FlagOrderFactory(goflect.UI_FIELDS),
	}
	appendTag := func(name, value string) {
		if _, hit := touchedTags[name]; !hit {
			if value != "" {
				if formatter, fieldHit := fieldFormatter[name]; fieldHit {
					value = formatter(value)
				}
				entries = append(entries, name+":"+value)
			}
			touchedTags[name] = 1
		}
	}
	for _, name := range goflect.TAGS {
		appendTag(name, tagKeys[name])
	}
	for name, value := range tagKeys {
		appendTag(name, value)
	}
	seperator := "\n\t"
	if cols > 0 {
		seperator += strings.Repeat(" ", cols)
	}
	return strings.Join(entries, seperator), errors
}

func FlagOrderFactory(flags []string) func(string) string {
	orderFlags := func(value string) string {
		wrapquotes, _ := regexp.Compile("(^\"|\"$)")
		commas, _ := regexp.Compile(", *")
		value = wrapquotes.ReplaceAllString(value, "")
		entries := commas.Split(value, -1)
		temp := make(map[string]int)
		for _, entry := range entries {
			temp[entry] = 1
		}
		output := make([]string, 0)
		for _, flag := range flags {
			if _, hit := temp[flag]; hit {
				output = append(output, flag)
			}
		}
		return strconv.Quote(strings.Join(output, ", "))

	}
	return orderFlags
}
