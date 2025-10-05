package testableexamples

type MyStruct struct {
	FirstName string
	LastName  string
}

func (m MyStruct) MyMethod() string {
	return m.FirstName + " " + m.LastName
}

func MyFunc(name string) string {
	return "Hello " + name + "!"
}
