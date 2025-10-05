package patterns

import (
	"log"
	"sync"
	"time"
)

// TODO works not as intended
func ExampleWorkerPool() {
	wp := NewWorkerPool[string](3)
	defer wp.Close()

	go func() {
		for res := range wp.Results() {
			log.Print(res.Res)
		}
	}()

	var wg sync.WaitGroup
	data := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	wg.Add(len(data))

	for _, dataItem := range data {
		di := dataItem
		wp.Go(func() (string, error) {
			time.Sleep(time.Second)
			wg.Done()
			return di, nil
		})
	}

	wg.Wait()
	// Output:
}
