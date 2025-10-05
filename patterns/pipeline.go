package patterns

func Pipeline(inputCh <-chan int) <-chan int {
	outputCh := make(chan int)
	go func() {
		// read from input
		for v := range inputCh {
			// do work
			// ...
			// send to output
			outputCh <- v
		}
		close(outputCh) // signal no more values are coming
	}()
	return outputCh
}
