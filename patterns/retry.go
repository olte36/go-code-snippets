package patterns

import (
	"context"
	"fmt"
	"time"
)

func Retry(fn func() error, maxRetries int, delay time.Duration) (int, error) {
	for i := range maxRetries {
		if err := fn(); err == nil {
			return i + 1, nil
		}
		<-time.After(delay)
	}
	return maxRetries, fmt.Errorf("max retries reached")
}

func RetryContext(ctx context.Context, fn func() error, maxRetries int, delay time.Duration) (int, error) {
	for i := range maxRetries {
		if err := fn(); err == nil {
			return i + 1, nil
		}
		select {
		case <-ctx.Done():
			return i + 1, ctx.Err()
		case <-time.After(delay):
			// continue
		}
	}
	return maxRetries, fmt.Errorf("max retries reached")
}
