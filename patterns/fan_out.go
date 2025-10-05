package patterns

import "context"

// Fan-out evenly distributes messages from an input channel to multiple output channels.
// It is a useful pattern for parallelizing CPU and I/O utilization.
func FanOut(in <-chan int, numChans int) []<-chan int {
	outs := make([]<-chan int, 0, numChans)

	// for every output channel read from the input channel in a separate goroutine
	for range numChans {
		ch := make(chan int)
		outs = append(outs, ch)
		go func(ch chan int) {
			defer close(ch)
			for v := range in {
				ch <- v
			}
		}(ch)
	}

	return outs
}

func FanOutContext(ctx context.Context, in <-chan int, numChans int) []<-chan int {
	outs := make([]<-chan int, 0, numChans)

	for range numChans {
		ch := make(chan int)
		outs = append(outs, ch)

		go func(ch chan int) {
			defer close(ch)
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					// write to the output channel
					select {
					case <-ctx.Done():
						return
					case ch <- v:
						// continue
					}
				}
			}
		}(ch)
	}

	return outs
}
