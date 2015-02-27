/*
Relational algrebra shows up in many forms in development.  However, there are rarely tools to support using the same relational constain with multiple contexts.  There are normal conditionals for in program logic, SQL for databases, and custom objects for performing stream evaluations.  Each of these have different serializtion methods, and different easy of use for composability.

This API 100% centered around the ability to compose relational expressions for reuse in different contexts.

The plan is to use matcher objects to define projections in the product.  This can be used in multiple ways

    A composable SQL querying API
    A filtering mechansim
    A way of validating data is valid for a struct
    A way of providing a filtering mechanism for user interaction
    A consistent rules API mechanism
*/
package matcher

type fieldOps int

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

type Matcher interface {
	Match(record interface{}) (bool, error)
}

type Yielder interface {
	Yield() (interface{}, error)
}

type Printer interface {
	Print(m Matcher) (string, error)
}
