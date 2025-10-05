package testableexamples

import (
	"fmt"
	"log"
)

func Example() {
	s := MyStruct{
		FirstName: "John",
		LastName:  "Dow",
	}
	_ = MyFunc(s.MyMethod())
}

func ExampleMyFunc() {
	res := MyFunc("Mike")
	fmt.Print(res)

	// Output: Hello Mike!
}

// this works because log prints to stderr
func ExampleMyFunc_logOutput() {
	res := MyFunc("Alice")
	log.Print(res)

	// Output:
}

func ExampleMyStruct() {
	s := MyStruct{
		FirstName: "John",
		LastName:  "Dow",
	}
	fmt.Printf("%+v", s)

	// Output: {FirstName:John LastName:Dow}
}

func ExampleMyStruct_MyMethod() {
	s := MyStruct{
		FirstName: "John",
		LastName:  "Dow",
	}
	res := s.MyMethod()
	fmt.Print(res)

	// Output: John Dow
}

func ExampleMyStruct_MyMethod_unorderedOutput() {
	s1 := MyStruct{
		FirstName: "John",
		LastName:  "Dow",
	}
	s2 := MyStruct{
		FirstName: "Sally",
		LastName:  "Gray",
	}
	for _, s := range []MyStruct{s2, s1} {
		fmt.Println(s.MyMethod())
	}

	// Unordered output: John Dow
	// Sally Gray
}
