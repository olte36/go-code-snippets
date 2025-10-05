package patterns

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	var curr atomic.Int32
	var wg sync.WaitGroup
	var semaphore = NewSemaphore(3)

	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()

			semaphore.Acquire()

			curr.Add(1)
			time.Sleep(3 * time.Second)
			fmt.Println("curr", curr.Load())
			curr.Add(-1)

			semaphore.Release()
		}()
	}
	wg.Wait()
}
