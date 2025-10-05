package patterns

import "context"

// Pipeline pattern is a way to process data through stages in a concurrent fashion.
func Pipeline(inputCh <-chan int, processFn func(int) int) <-chan int {
	outputCh := make(chan int)
	go func() {
		defer close(outputCh) // signal no more values are coming

		// read from input
		for v := range inputCh {
			// do work
			r := processFn(v)
			// send to output
			outputCh <- r
		}
	}()
	return outputCh
}

func PipelineContext(ctx context.Context, inputCh <-chan int, processFn func(int) int) <-chan int {
	outputCh := make(chan int)
	go func() {
		defer close(outputCh)

		for {
			// read from input
			select {
			case <-ctx.Done():
				return
			case v, ok := <-inputCh:
				if !ok {
					return
				}
				// process the value
				r := processFn(v)
				// send to output
				select {
				case <-ctx.Done():
					return
				case outputCh <- r:
					// continue
				}
			}
		}
	}()
	return outputCh
}
