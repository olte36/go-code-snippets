package patterns

import (
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
