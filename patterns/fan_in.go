package patterns

import "sync"

func FanIn(chans []chan int) chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch chan int) {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
