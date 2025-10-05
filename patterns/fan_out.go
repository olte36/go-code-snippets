package patterns

func FanOut(in chan int, numChans int) []chan int {
	outs := make([]chan int, numChans)

	for i := range numChans {
		outs[i] = make(chan int)
		go func() {
			defer close(outs[i])
			for v := range in {
				outs[i] <- v
			}
		}()
	}

	return outs
}
