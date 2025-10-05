package patterns

import (
	"context"
	"log"
	"time"
)

func Retry(fn TargetFunc, maxRetries int, delay time.Duration) TargetFunc {
	return func(ctx context.Context) error {
		for i := 0; ; i++ {
			if err := fn(ctx); err == nil || i >= maxRetries {
				return err
			}

			log.Printf("Attempt number %d failed, retrying in %v", i+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
