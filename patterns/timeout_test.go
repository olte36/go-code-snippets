package patterns

import (
	"fmt"
	"go.uber.org/goleak"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	defer goleak.VerifyNone(t)

	_, err := Timeout(func() int {
		time.Sleep(2 * time.Second)
		return 1
	}, 1*time.Second)

	time.Sleep(2 * time.Second)

	fmt.Println(err)
}
