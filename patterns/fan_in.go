package patterns

import (
	"context"
	"sync"
)

// Fan-in multiplexes multiple input channels onto one output channel.
// Services that have some number of workers that all generate output may find it useful
// to combine all of the workers’ outputs to be processed as a single unified stream.
func FanIn(chans []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(chans))
	// read from input channels in separate goroutines
	// and send the values to the output channel
	for _, ch := range chans {
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}(ch)
	}
	// close the out channel when all reading is done
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func FanInContext(ctx context.Context, chans []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch <-chan int) {
			defer wg.Done()
			for {
				// read from the input channel
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}
					// write to the output channel
					select {
					case <-ctx.Done():
						return
					case out <- v:
						// continue
					}
				}
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
