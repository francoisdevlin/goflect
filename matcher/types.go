/*
Relational algrebra shows up in many forms in development.  However, there are rarely tools to support using the same relational constain with multiple contexts.  There are normal conditionals for in program logic, SQL for databases, and custom objects for performing stream evaluations.  Each of these have different serializtion methods, and different easy of use for composability.

This API 100% centered around the ability to compose relational expressions for reuse in different contexts.

The plan is to use matcher objects to define projections in the product.  This can be used in multiple ways

    A composable SQL querying API
    A filtering mechansim
    A way of validating data is valid for a struct
    A way of providing a filtering mechanism for user interaction
    A consistent rules API mechanism
    A way of determining set membership, for projects like HOT
*/
package matcher

type fieldOps int

/*
This is the complete list of allowed field ops
*/
const (
	LT fieldOps = iota
	LTE
	GT
	GTE
	EQ
	NEQ
	IN
	NOT_IN
	MATCH
	NOT_MATCH
)

/*
This is the core abstraction in the matcher library. It is designed to be a conditional that can be passed around.
*/
type Matcher interface {
	Match(record interface{}) (bool, error)
}

/*
This is a specialization of the matcher designed to match a struct or dictionary.  It is designed to interact with the outside world as such

	* Filter a stream of structs
	* Serialize as a SQL WHERE clause
	* Compose multiple SQL views on the fly
	* Interact with the end user by using the pretty printer and parsing routines
*/
type StructMatcher interface {
	Matcher
	AddField(string, Matcher)
	Field(string) Yielder
}

/*
This is an interface for parsing a string, and creating a matcher from it.  It requires a context, which is provided to the constructor (see NewParser for more information about this).  It uses the context to determine what symbols are valid, and if it is possible to use the type in the resultant matcher
*/
type Parser interface {
	Parse(string) (Matcher, error)
}

/*
The yielder is used to provide dynmaic values in a matcher.  It is needs let expressions such as 'A = B' work properly.  It can also be used to determine how an item compares to a dynamic value, such as a time window (e.g. "Last 5 minutes")
*/
type Yielder interface {
	Yield() (interface{}, error)
}

type Printer interface {
	Print(m Matcher) (string, error)
}
