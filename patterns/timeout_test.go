package patterns

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/goleak"
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

func TestTimeoutContext(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := TimeoutContext(ctx, func() int {
		time.Sleep(2 * time.Second)
		return 1
	}, 1*time.Second)

	time.Sleep(2 * time.Second)

	fmt.Println(err)
}
