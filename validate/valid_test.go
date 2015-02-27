package validate

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/goflect"
)

func ExampleDefaultMatcher_1() {
	type PasswordChange struct {
		Current string
		New     string `valid:"New != Current AND New == Repeat"`
		Repeat  string
	}
	fields := goflect.GetInfo(PasswordChange{})
	def, _ := DefaultMatcher(fields[0])
	fmt.Println("Test")
	def.Match(nil)

	//Output:
	//Bacon
}
