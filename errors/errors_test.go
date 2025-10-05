package errors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ConstErr = &PointerImplErr{text: "very bad thing happened"}

// TestErrorsIsAnsAsDemo demonstrates and tests the behavior of Go's errors.Is and errors.As functions
// with various error types and scenarios. It covers cases including value and pointer errors, custom
// Is implementations, and correct and incorrect usage of errors.As with different target types.
// The test verifies expected outcomes and panics for invalid usages, providing insight into how
// error comparison and type assertion work in Go's error handling patterns.
func TestErrorsIsAnsAsDemo(t *testing.T) {
	testCases := []struct {
		name     string
		fn       func() bool
		expected bool
		panics   bool
	}{
		{
			name: "errors.Is with value error",
			// the underlying data is a struct, the structs are compared
			fn: func() bool {
				var err, target error

				err = fmt.Errorf("%w", ValueImplErr{})
				target = ValueImplErr{}

				return errors.Is(err, target)
			},
			expected: true,
		},
		{
			name: "errors.Is with pointer error",
			// the underlying data is a pointer, the pointers are compared,
			// they point to different data
			fn: func() bool {
				var err, target error

				err = fmt.Errorf("%w", &PointerImplErr{})
				target = &PointerImplErr{}

				return errors.Is(err, target)
			},
			expected: false,
		},
		{
			name: "errors.Is with pointer error, same data",
			// the underlying data is a pointer, the pointers are compared,
			// they point to the same data
			fn: func() bool {
				var err, target error

				err = fmt.Errorf("%w", ConstErr)
				target = ConstErr

				return errors.Is(err, target)
			},
			expected: true,
		},
		{
			name: "errors.Is with custom Is implementation",
			fn: func() bool {
				var err, target error

				err = fmt.Errorf("%w", &CustomIsImplErr{})
				target = &CustomIsImplErr{}
				// true
				// CustomCompErr implements Is method
				return errors.Is(err, target)
			},
			expected: true,
		},
		{
			name: "errors.As with value error",
			// target in errors.As must be a pointer to the data of the type
			// that implements error interface
			// In this case it is a pointer to the struct
			fn: func() bool {
				var err error
				var target *ValueImplErr

				err = fmt.Errorf("%w", ValueImplErr{})
				target = &ValueImplErr{}

				return errors.As(err, target)
			},
			expected: true,
		},
		{
			name: "errors.As with pointer error",
			// target in errors.As must be a pointer to the data of the type
			// that implements error interface
			// In this case it is a pointer to the pointer
			fn: func() bool {
				var err error
				var target **PointerImplErr

				err = fmt.Errorf("%w", &PointerImplErr{})
				x := &PointerImplErr{}
				target = &x

				return errors.As(err, target)
			},
			expected: true,
		},
		{
			name: "errors.As with pointer error, wrong target type",
			// target in errors.As must be a pointer to the data of the type
			// that implements error interface
			// PointerImplErr does not implement error interface
			// but *PointerImplErr does, so target must be **PointerImplErr
			fn: func() bool {
				var err error
				var target any // *PointerImplErr

				err = fmt.Errorf("%w", &PointerImplErr{})
				target = &PointerImplErr{}

				return errors.As(err, target)
			},
			panics: true,
		},
		{
			name: "errors.As with interface target implemented by value",
			fn: func() bool {
				var err error
				var target any // DetailsInterface

				err = fmt.Errorf("%w", DetailsValueImplErr{})
				var x DetailsInterface = DetailsValueImplErr{}
				target = x

				return errors.As(err, target)
			},
			panics: true,
		},
		{
			name: "errors.As with pointer to interface target implemented by value",
			fn: func() bool {
				var err error
				var target *DetailsInterface

				err = fmt.Errorf("%w", DetailsValueImplErr{})
				var x DetailsInterface = DetailsValueImplErr{}
				target = &x

				return errors.As(err, target)
			},
			expected: true,
		},
		{
			name: "errors.As with interface target implemented by pointer",
			fn: func() bool {
				var err error
				var target any // DetailsInterface

				err = fmt.Errorf("%w", &DetailsPointerImplErr{})
				var x DetailsInterface = &DetailsPointerImplErr{}
				target = x

				return errors.As(err, target)
			},
			panics: true,
		},
		{
			name: "errors.As with pointer to interface target implemented by pointer",
			fn: func() bool {
				var err error
				var target *DetailsInterface

				err = fmt.Errorf("%w", &DetailsPointerImplErr{details: "bad thing happened"})
				var x DetailsInterface = &DetailsPointerImplErr{}
				target = &x

				return errors.As(err, target)
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics {
				assert.Panics(t, func() {
					_ = tc.fn()
				})
			} else {
				assert.Equal(t, tc.expected, tc.fn())
			}
		})
	}
}

var _ error = ValueImplErr{}

type ValueImplErr struct {
	text string
}

func (v ValueImplErr) Error() string {
	return "SimpleValueErr: " + v.text
}

var _ error = (*PointerImplErr)(nil)

type PointerImplErr struct {
	text string
}

func (p *PointerImplErr) Error() string {
	return "PointerImplErr: " + p.text
}

type DetailsInterface interface {
	Details() string
}

var _ error = DetailsValueImplErr{}
var _ DetailsInterface = DetailsValueImplErr{}

type DetailsValueImplErr struct {
	details string
}

func (d DetailsValueImplErr) Error() string {
	return "DetailsValueErr" + d.details
}

func (d DetailsValueImplErr) Details() string {
	return d.details
}

var _ error = (*DetailsPointerImplErr)(nil)
var _ DetailsInterface = (*DetailsPointerImplErr)(nil)

type DetailsPointerImplErr struct {
	details string
}

func (d *DetailsPointerImplErr) Error() string {
	return "DetailsPointerErr" + d.details
}

func (d *DetailsPointerImplErr) Details() string {
	return d.details
}

var _ error = (*CustomIsImplErr)(nil)
var _ (interface{ Is(error) bool }) = (*CustomIsImplErr)(nil)

type CustomIsImplErr struct {
	text string
}

func (c *CustomIsImplErr) Error() string {
	return "CustomIsImplErr: " + c.text
}

func (c *CustomIsImplErr) Is(target error) bool {
	t, ok := target.(*CustomIsImplErr)
	return ok && c.text == t.text
}
