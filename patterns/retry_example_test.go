package patterns

import (
	"context"
	"errors"
	"log"
	"time"
)

func ExampleRetry() {
	fn := failCountTimes(2)

	fn = Retry(fn, 4, 500*time.Millisecond)

	if err := fn(context.Background()); err != nil {
		log.Print(err)
	}
	// Output:
}

func ExampleRetry_ctxCanceled() {
	fn := failCountTimes(3)

	fn = Retry(fn, 4, 1*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	if err := fn(ctx); err != nil {
		log.Print(err)
	}
	// Output:
}

func failCountTimes(count int) TargetFunc {
	var c int
	return func(context.Context) error {
		c++
		if c <= count {
			return errors.New("service unavailable")
		}
		log.Print("Successfull call")
		return nil
	}
}
