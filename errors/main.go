package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var VeryBadThingHappenedErr = &SimplePointerErr{text: "very bad thing happened"}

func main() {
	{
		var err, target error
		err = fmt.Errorf("%w", SimpleValueErr{text: "bad thing happened"})
		target = SimpleValueErr{text: "bad thing happened"}
		// true
		// the underlying data is a struct, the structs are compared
		fmt.Println("1)", errors.Is(err, target))
	}
	{
		var err, target error
		err = fmt.Errorf("%w", &SimplePointerErr{text: "bad thing happened"})
		target = &SimplePointerErr{text: "bad thing happened"}
		// false
		// the underlying data is a pointer, the pointers are compared,
		// they point to different data
		fmt.Println("2)", errors.Is(err, target))
	}
	{
		var err, target error
		err = fmt.Errorf("%w", VeryBadThingHappenedErr)
		target = VeryBadThingHappenedErr
		// true
		// the underlying data is a pointer, the pointers are compared,
		// they point to the same data
		fmt.Println("3)", errors.Is(err, target))
	}
	{
		var err error
		var target *SimpleValueErr
		err = fmt.Errorf("%w", SimpleValueErr{text: "bad thing happened"})
		target = &SimpleValueErr{}
		// true
		// target in errors.As must be a pointer to the data of the type
		// that implements error interface
		// In this case it is a pointer to the struct
		fmt.Println("4)", errors.As(err, target), target.Error())
	}
	{
		var err error
		var target *SimplePointerErr
		err = fmt.Errorf("%w", &SimplePointerErr{text: "bad thing happened"})
		target = &SimplePointerErr{}
		func() {
			defer func() {
				fmt.Println("5)", recover())
			}()
			// panics
			// target in errors.As must be a pointer to the data of the type
			// that implements error interface
			// SimplePointerErr does not implement error interface
			// but *SimplePointerErr, so target must be **SimplePointerErr
			errors.As(err, target)
		}()
	}
	{
		var err error
		var x *SimplePointerErr
		var target **SimplePointerErr
		err = fmt.Errorf("%w", &SimplePointerErr{text: "bad thing happened"})
		x = &SimplePointerErr{}
		target = &x
		// true
		// target in errors.As must be a pointer to the data of the type
		// that implements error interface
		// In this case it is a pointer to the pointer
		fmt.Println("6)", errors.As(err, target), x.Error())
	}
	{
		var err error
		var target DetailsInterface
		err = fmt.Errorf("%w", DetailsValueErr{details: "bad thing happened"})
		target = &DetailsValueErr{}
		// true
		fmt.Println("7)", errors.As(err, target), target.Details())
	}
	{
		var err error
		var x DetailsInterface
		var target *DetailsInterface
		err = fmt.Errorf("%w", DetailsValueErr{details: "bad thing happened"})
		x = DetailsValueErr{}
		target = &x
		// true
		fmt.Println("8)", errors.As(err, target), x.Details())
	}
	{
		var err error
		var target DetailsInterface
		err = fmt.Errorf("%w", &DetailsPointerErr{details: "bad thing happened"})
		target = &DetailsPointerErr{}
		func() {
			defer func() {
				fmt.Println("9)", recover())
			}()
			errors.As(err, target)
		}()
	}
	{
		var err error
		var x DetailsInterface
		var target *DetailsInterface
		err = fmt.Errorf("%w", &DetailsPointerErr{details: "bad thing happened"})
		x = &DetailsPointerErr{}
		target = &x
		fmt.Println("10)", errors.As(err, target), x.Details())
	}
	{
		var err error
		var target *CustomCompErr
		err = fmt.Errorf("%w", &CustomCompErr{details: []string{"1", "2"}})
		target = &CustomCompErr{details: []string{"1", "2"}}
		// true
		// CustomCompErr implements Is method
		fmt.Println("11)", errors.Is(err, target))
	}
}

var _ error = SimpleValueErr{}
var _ error = (*SimplePointerErr)(nil)
var _ error = DetailsValueErr{}
var _ DetailsInterface = DetailsValueErr{}
var _ error = (*DetailsPointerErr)(nil)
var _ DetailsInterface = (*DetailsPointerErr)(nil)
var _ error = (*CustomCompErr)(nil)
var _ (interface{ Is(error) bool }) = (*CustomCompErr)(nil)
var _ (interface{ As(any) bool }) = (*CustomCompErr)(nil)

type SimpleValueErr struct {
	text string
}

type SimplePointerErr struct {
	text string
}

type DetailsInterface interface {
	Details() string
}

type DetailsValueErr struct {
	details string
}

type DetailsPointerErr struct {
	details string
}

type CustomCompErr struct {
	details []string
}

func (s SimpleValueErr) Error() string {
	return "SimpleValueErr: " + s.text
}

func (s *SimplePointerErr) Error() string {
	return "SimplePointerErr: " + s.text
}

func (d DetailsValueErr) Error() string {
	return "DetailsValueErr" + d.details
}

func (d DetailsValueErr) Details() string {
	return d.details
}

func (d *DetailsPointerErr) Error() string {
	return "DetailsPointerErr" + d.details
}

func (d *DetailsPointerErr) Details() string {
	return d.details
}

func (c *CustomCompErr) Error() string {
	return "CustomCompErr: " + strings.Join(c.details, ", ")
}

func (c *CustomCompErr) Is(target error) bool {
	t, ok := target.(*CustomCompErr)
	return ok && slices.Equal(c.details, t.details)
}

func (c *CustomCompErr) As(target any) bool {
	return errors.As(c, target)
}
