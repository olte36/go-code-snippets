package patterns

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRetry_Successfull(t *testing.T) {
	fn := failCountTimes(2)

	retried, err := Retry(fn, 4, 500*time.Millisecond)

	fmt.Println("Retried:", retried, "times, error", err)
}

func TestRetry_Failed(t *testing.T) {
	fn := failCountTimes(4)

	retried, err := Retry(fn, 2, 500*time.Millisecond)

	fmt.Println("Retried:", retried, "times, error:", err)
}

func failCountTimes(count int) func() error {
	var c int
	return func() error {
		c++
		if c <= count {
			return errors.New("service unavailable")
		}
		fmt.Println("Successfull call")
		return nil
	}
}
