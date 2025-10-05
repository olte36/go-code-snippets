package patterns

import (
	"context"
	"fmt"
	"time"
)

func Timeout(fn func() int, timeout time.Duration) (int, error) {
	resCh := make(chan int)
	go func() {
		select {
		case resCh <- fn():
		default: // prevent goroutine leak
		}
	}()
	select {
	case <-time.After(timeout):
		return 0, fmt.Errorf("timed out")
	case res := <-resCh:
		return res, nil
	}
}

func TimeoutContext(cxt context.Context, fn func() int, timeout time.Duration) (int, error) {
	resCh := make(chan int)

	go func() {
		select {
		case resCh <- fn():
		case <-cxt.Done():
		}
	}()

	select {
	case <-cxt.Done():
		return 0, cxt.Err()
	case <-time.After(timeout):
		return 0, fmt.Errorf("timed out")
	case res := <-resCh:
		return res, nil
	}
}
