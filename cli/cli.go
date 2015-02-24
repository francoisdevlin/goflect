/*
This is a tool to marshal command line flags into a concrete type.  See the examples of the FlagSetup function to learn how it works
*/
package goflect

import (
	"flag"
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"reflect"
	"strconv"
)

/*
This function expects to recieve a set of arguments from the command line, and use them to hydrate a struct.  It is a simple wrapper around Go's flag package, and leverages much of that tools functionality

The tool makes extensive use of the Default and Description metadata
*/
func FlagSetup(record interface{}, args []string) {
	fields := goflect.GetInfo(record)
	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	vals := make(map[string]interface{})
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	for _, field := range fields {
		switch field.Kind {
		case reflect.Bool:
			b, _ := strconv.ParseBool(field.Default)
			vals[field.Name] = flagSet.Bool(field.Name, b, field.Description)
		case reflect.Float64, reflect.Float32:
			f, _ := strconv.ParseFloat(field.Default, 64)
			vals[field.Name] = flagSet.Float64(field.Name, f, field.Description)
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			i, _ := strconv.ParseInt(field.Default, 10, 64)
			vals[field.Name] = flagSet.Int64(field.Name, i, field.Description)
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			u, _ := strconv.ParseUint(field.Default, 10, 64)
			vals[field.Name] = flagSet.Uint64(field.Name, u, field.Description)
		default:
			vals[field.Name] = flagSet.String(field.Name, field.Default, field.Description)
		}
	}
	flagSet.Parse(args[1:])
	for _, field := range fields {
		fieldVal := val.FieldByName(field.Name)
		var temp reflect.Value
		switch field.Kind {
		case reflect.Bool:
			temp = (reflect.ValueOf(*vals[field.Name].(*bool)))
		case reflect.Float64:
			temp = (reflect.ValueOf(float64(*vals[field.Name].(*float64))))
		case reflect.Float32:
			temp = (reflect.ValueOf(float32(*vals[field.Name].(*float64))))
		case reflect.Int:
			temp = (reflect.ValueOf(int(*vals[field.Name].(*int64))))
		case reflect.Int64:
			temp = (reflect.ValueOf(int64(*vals[field.Name].(*int64))))
		case reflect.Int32:
			temp = (reflect.ValueOf(int32(*vals[field.Name].(*int64))))
		case reflect.Int16:
			temp = (reflect.ValueOf(int16(*vals[field.Name].(*int64))))
		case reflect.Int8:
			temp = (reflect.ValueOf(int8(*vals[field.Name].(*int64))))
		case reflect.Uint:
			temp = (reflect.ValueOf(uint(*vals[field.Name].(*uint64))))
		case reflect.Uint64:
			temp = (reflect.ValueOf(uint64(*vals[field.Name].(*uint64))))
		case reflect.Uint32:
			temp = (reflect.ValueOf(uint32(*vals[field.Name].(*uint64))))
		case reflect.Uint16:
			temp = (reflect.ValueOf(uint16(*vals[field.Name].(*uint64))))
		case reflect.Uint8:
			temp = (reflect.ValueOf(uint8(*vals[field.Name].(*uint64))))
		default:
			temp = (reflect.ValueOf(*vals[field.Name].(*string)))
		}
		fieldVal.Set(temp)
	}

}
