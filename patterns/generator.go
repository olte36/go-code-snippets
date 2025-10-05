package patterns

func Generator(count int) <-chan int {
	ch := make(chan int)
	// generator should not block, so generate values in a goroutine
	go func() {
		// do work
		for i := range count {
			ch <- i + 1
		}
		close(ch) // signal no more values are coming
	}()
	return ch
}
