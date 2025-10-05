package patterns

import "sync"

func Tee(input <-chan int, numChs int) []<-chan int {
	outChs := make([]chan int, numChs)
	for i := range numChs {
		outChs[i] = make(chan int)
	}

	go func() {
		defer func() {
			for _, ch := range outChs {
				close(ch)
			}
		}()

		var wg sync.WaitGroup
		for v := range input {
			wg.Add(numChs)
			for i := range numChs {
				go func(<-chan int) {
					defer wg.Done()
					outChs[i] <- v
				}(outChs[i])
			}
			wg.Wait()
		}
	}()

	// remap to read-only channels
	resChs := make([]<-chan int, numChs)
	for i := range numChs {
		resChs[i] = outChs[i]
	}
	return resChs
}
