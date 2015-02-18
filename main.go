package main

import (
	//"database/sql"
	//"strings"
	//"html/template"
	//"net/http"
	//"flag"
	"fmt"
	"os"
	//"reflect"
	//"strconv"
)

func main() {
	type Bar struct {
		I       int `desc:"This is the Id"`
		I64     int64
		I32     int32
		I16     int16
		I8      int8
		U       uint `desc:"This is the Ud" default:"1234"`
		U64     uint64
		U32     uint32
		U16     uint16
		U8      uint8
		F64     float64
		F32     float32
		S       string
		Awesome bool `desc:"This controls the awesomeness"`
	}

	//record := Info{}
	record := Bar{}
	FlagSetup(&record, os.Args)
	//fields := GetInfo(&record)

	//flags := flag.FlagSet{}

	//for _, field := range fields {
	////fmt.Println(field)
	//flag.String(field.Name, "", field.Description)
	//flags = flags
	//}

	fmt.Println(record)

	//handler := func(w http.ResponseWriter, r *http.Request) {
	////t, _ := template.ParseFiles("header.html")
	////t.Execute(w, "My Title")
	//t, _ := template.ParseFiles("test.html")
	//t.Execute(w, records)
	////fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	////fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	//}
	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":8080", nil)
}
