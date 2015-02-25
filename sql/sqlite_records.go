package records

//import (
//"database/sql"
//"fmt"
//"git.sevone.com/sdevlin/goflect.git/goflect"
//"reflect"
//"strconv"
//"strings"
//)

//func typeAndVal(record interface{}) (reflect.Type, reflect.Value) {
//typ := reflect.TypeOf(record)
//// if a pointer to a struct is passed, get the type of the dereferenced object
//if typ.Kind() == reflect.Ptr {
//typ = typ.Elem()
//}

//val := reflect.ValueOf(record)
//if val.Kind() == reflect.Ptr {
//val = val.Elem()
//}
//return typ, val
//}

//func sqliteLookupMap() map[reflect.Kind]string {
//lookup := map[reflect.Kind]string{
//reflect.Bool:    "integer",
//reflect.Int:     "integer",
//reflect.Int8:    "integer",
//reflect.Int16:   "integer",
//reflect.Int32:   "integer",
//reflect.Int64:   "integer",
//reflect.Uint:    "integer",
//reflect.Uint8:   "integer",
//reflect.Uint16:  "integer",
//reflect.Uint32:  "integer",
//reflect.Uint64:  "integer",
//reflect.Float32: "real",
//reflect.Float64: "real",
//}
//return lookup
//}

//type SqliteRecordService struct {
//Conn *sql.DB
//}

//func (service SqliteRecordService) Insert(record interface{}) {
//typ, val := typeAndVal(record)

//fields := goflect.GetInfo(record)
//statement := ""
//statement += "INSERT INTO " + typ.Name() + "("
//columns := make([]string, 0)
//for _, field := range fields {
//if field.IsAutoincrement {
//continue
//}
//columns = append(columns, field.Name)
//}
//statement += strings.Join(columns, ", ")
//statement += " ) VALUES ("

//columns = make([]string, 0)
//for _, field := range fields {
//if field.IsAutoincrement {
//continue
//}
//fieldVal := val.FieldByName(field.Name)
//columns = append(columns, wrap(fieldVal, field))
//}
//statement += strings.Join(columns, ", ")
//statement += " )"
//_, _ = service.Conn.Exec(statement)
//}

//func (service SqliteRecordService) Update(record interface{}) {
//typ, val := typeAndVal(record)

//fields := goflect.GetInfo(record)
//statement := ""
//statement += "UPDATE " + typ.Name() + " SET "
//conditions := make(map[string]interface{})
//columns := make([]string, 0)
//for _, field := range fields {
//fieldVal := val.FieldByName(field.Name)
//if field.IsPrimary {
//switch fieldVal.Kind() {
//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//conditions[field.Name] = fieldVal.Uint()
//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//conditions[field.Name] = fieldVal.Int()
//default:
//conditions[field.Name] = fieldVal.Int()
//continue
//}
//}
//temp := field.Name + "=" + wrap(fieldVal, field)
//columns = append(columns, temp)
//}
//statement += strings.Join(columns, ", ")

//statement += processWhereClause(fields, conditions)
////fmt.Println(statement)
//service.Conn.Exec(statement)
//}

//func (service SqliteRecordService) Delete(record interface{}) {
//typ, val := typeAndVal(record)

//fields := goflect.GetInfo(record)
//statement := ""
//statement += "DELETE FROM " + typ.Name()
//conditions := make(map[string]interface{})
//for _, field := range fields {
//fieldVal := val.FieldByName(field.Name)
//if field.IsPrimary {
//switch fieldVal.Kind() {
//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//conditions[field.Name] = fieldVal.Uint()
//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//conditions[field.Name] = fieldVal.Int()
//default:
//conditions[field.Name] = fieldVal.Int()
//continue
//}
//}
//}

//statement += processWhereClause(fields, conditions)
//service.Conn.Exec(statement)
//}

//func (service SqliteRecordService) DeleteById(id int, record interface{}) {
//typ, _ := typeAndVal(record)

//fields := goflect.GetInfo(record)
//statement := ""
//statement += "DELETE FROM " + typ.Name()
//conditions := make(map[string]interface{})
//for _, field := range fields {
//if field.IsPrimary {
//conditions[field.Name] = id
//}
//}

//statement += processWhereClause(fields, conditions)
//service.Conn.Exec(statement)
//}

//func (service SqliteRecordService) ReadAll(record interface{}) func(record interface{}) bool {
//conditions := make(map[string]interface{})
//return service.ReadAllWhere(record, conditions)
//}

//func (service SqliteRecordService) Get(id int64, record interface{}) {
//fields := goflect.GetInfo(record)
//conditions := make(map[string]interface{})
//for _, field := range fields {
//if field.IsPrimary {
//conditions[field.Name] = id
//}
//}
//next := service.ReadAllWhere(record, conditions)
//for next(record) {
//} //The last call closes the result set, important!
//}

//func (service SqliteRecordService) GetByNominal(name string, record interface{}) {
//fields := goflect.GetInfo(record)
//conditions := make(map[string]interface{})
//for _, field := range fields {
//if field.IsNominal {
//conditions[field.Name] = name
//}
//}
//next := service.ReadAllWhere(record, conditions)
//for next(record) {
//} //The last call closes the result set, important!
//}

//func (service SqliteRecordService) ReadAllWhere(record interface{}, conditions map[string]interface{}) func(record interface{}) bool {
//typ, _ := typeAndVal(record)

//fields := goflect.GetInfo(record)
//statement := "SELECT "
//columns := make([]string, 0)
//for _, field := range fields {
//columns = append(columns, field.Name)
//}
//statement += strings.Join(columns, " , ")
//statement += " FROM " + typ.Name()

//statement += processWhereClause(fields, conditions)

//rows, _ := service.Conn.Query(statement)

//output := func(r interface{}) bool {
//return nextRow(rows, r)
//}

//return output
//}

//func (service SqliteRecordService) ReadAllNominal(record interface{}) func(record *goflect.Nominal) bool {
//conditions := make(map[string]interface{})
//return service.ReadAllNominalWhere(record, conditions)
//}

//func (service SqliteRecordService) ReadAllNominalWhere(record interface{}, conditions map[string]interface{}) func(record *goflect.Nominal) bool {
//typ, _ := typeAndVal(record)

//fields := goflect.GetInfo(record)
//statement := "SELECT "
//columns := make([]string, 0)
//for _, field := range fields {
//if field.IsPrimary || field.IsNominal {
//columns = append(columns, field.Name)
//}
//}
//statement += strings.Join(columns, " , ")
//statement += " FROM " + typ.Name()

//statement += processWhereClause(fields, conditions)

//rows, _ := service.Conn.Query(statement)

//output := func(r *goflect.Nominal) bool {
//return nextRow(rows, r)
//}

//return output
//}

//func CreateSQLiteTable(record interface{}) (statement string) {
//typ, _ := typeAndVal(record)

//fields := goflect.GetInfo(record)
//lookup := sqliteLookupMap()
//statement = ""
//statement += "CREATE TABLE IF NOT EXISTS " + typ.Name() + "("
//for i, field := range fields {
//kind, present := lookup[field.Kind]
//if !present {
//kind = "string"
//}

//statement += "\n" + field.Name + " " + kind
//if field.IsPrimary {
//statement += " primary key"
//}
//if field.IsAutoincrement {
//statement += " autoincrement"
//}
//if !field.IsNullable {
//statement += " not null"
//}
//if i != len(fields)-1 {
//statement += ","
//}
//}
//statement += "\n)"
//return statement
//}

//func wrap(fieldVal reflect.Value, field goflect.Info) string {
//output := ""
//switch fieldVal.Kind() {
//case reflect.Bool:
//if fieldVal.Bool() {
//output = "1"
//} else {
//output = "0"
//}
////output = "" + strconv.FormatBool(fieldVal.Bool()) + ""
//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//output = "" + strconv.FormatInt(int64(fieldVal.Uint()), 10) + ""
//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//output = "" + strconv.FormatInt(fieldVal.Int(), 10) + ""
//case reflect.Float32:
//output = "" + strconv.FormatFloat(fieldVal.Float(), 'f', 10, 32) + ""
//case reflect.Float64:
//output = "" + strconv.FormatFloat(fieldVal.Float(), 'f', 10, 64) + ""
////break
//default:
//output = "\"" + fieldVal.String() + "\""
//}
//return output
//}

//func processWhereClause(fields []goflect.Info, conditions map[string]interface{}) string {
//output := ""
//segments := make([]string, 0)
//for _, field := range fields {
//if conditional, present := conditions[field.Name]; present {
//condVal := reflect.ValueOf(conditional)
//if condVal.Kind() == reflect.Ptr {
//condVal = condVal.Elem()
//}
//segments = append(segments, fmt.Sprintf("%v = %v", field.Name, wrap(condVal, field)))
//}
//}
//if len(segments) > 0 {
//output += " WHERE " + strings.Join(segments, " AND ")
//}
//return output
//}

//func nextRow(rows *sql.Rows, record interface{}) bool {
//next := rows.Next()
//if next {
//fields := goflect.GetInfo(record)
//val := reflect.ValueOf(record)
//if val.Kind() == reflect.Ptr {
//val = val.Elem()
//}

//vals := make([]interface{}, len(fields))
//addrs := make([]interface{}, len(fields))
//for i, _ := range vals {
//addrs[i] = &vals[i]
//}
//rows.Scan(addrs...)
//for i, field := range fields {
//fieldVal := val.FieldByName(field.Name)
//switch field.Kind {
//case reflect.Bool:
//fieldVal.Set(reflect.ValueOf(vals[i].(int64) != 0))
//case reflect.Float64:
//fieldVal.Set(reflect.ValueOf(float64(vals[i].(float64))))
//case reflect.Float32:
//fieldVal.Set(reflect.ValueOf(float32(vals[i].(float64))))
//case reflect.Int:
//fieldVal.Set(reflect.ValueOf(int(vals[i].(int64))))
//case reflect.Int64:
//fieldVal.Set(reflect.ValueOf(vals[i]))
//case reflect.Int32:
//fieldVal.Set(reflect.ValueOf(int32(vals[i].(int64))))
//case reflect.Int16:
//fieldVal.Set(reflect.ValueOf(int16(vals[i].(int64))))
//case reflect.Int8:
//fieldVal.Set(reflect.ValueOf(int8(vals[i].(int64))))
//case reflect.Uint:
//fieldVal.Set(reflect.ValueOf(uint(vals[i].(int64))))
//case reflect.Uint64:
//fieldVal.Set(reflect.ValueOf(uint64(vals[i].(int64))))
//case reflect.Uint32:
//fieldVal.Set(reflect.ValueOf(uint32(vals[i].(int64))))
//case reflect.Uint16:
//fieldVal.Set(reflect.ValueOf(uint16(vals[i].(int64))))
//case reflect.Uint8:
//fieldVal.Set(reflect.ValueOf(uint8(vals[i].(int64))))
//default:
//fieldVal.Set(reflect.ValueOf(string(vals[i].([]uint8))))
//}
//}

//}
//return next
//}
