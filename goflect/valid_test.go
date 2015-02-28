package goflect

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
Here we can see how to specify the default matcher for the password change struct. This allows us to define stricter constaints than type types alone would allow
*/
func ExampleDefaultMatcher_1() {
	type Password struct {
		Current string
		New     string `valid:"New != Current AND New = Repeat"`
		Repeat  string
	}

	def, err := DefaultMatcher(&Password{})
	if err != nil {
		fmt.Println(err)
		return
	}
	storedPassword := "old secret"
	old := matcher.NewStructMatcher()
	old.AddField("Current", matcher.Eq(storedPassword))
	m := matcher.And(def, old)

	passwordRequest := Password{
		Current: "old secret",
		New:     "new secret",
		Repeat:  "new secret",
	}
	if match, _ := m.Match(passwordRequest); match {
		fmt.Println("The password request is well formed")
	}

	passwordRequest = Password{
		Current: "mistake",
		New:     "new secret",
		Repeat:  "new secret",
	}
	if match, _ := m.Match(passwordRequest); !match {
		fmt.Println("The password request failed, the current is wrong")
	}

	passwordRequest = Password{
		Current: "old secreet",
		New:     "new secret",
		Repeat:  "wrong password repeated",
	}
	if match, _ := m.Match(passwordRequest); !match {
		fmt.Println("The password request failed, the wrong password repeated")
	}

	passwordRequest = Password{
		Current: "old secreet",
		New:     "old secret",
		Repeat:  "old secret",
	}
	if match, _ := m.Match(passwordRequest); !match {
		fmt.Println("The password request failed, the secret was not changed")
	}

	//Output:
	//The password request is well formed
	//The password request failed, the current is wrong
	//The password request failed, the wrong password repeated
	//The password request failed, the secret was not changed
}
