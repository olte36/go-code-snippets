package patterns

import (
	"fmt"
	"go.uber.org/goleak"
	"sync"
	"testing"
	"time"
)

func TestTee(t *testing.T) {
	defer goleak.VerifyNone(t)

	chans := Tee(Generator(5), 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range chans[0] {
			time.Sleep(1 * time.Second)
			fmt.Println("Reader 1:", v)
		}
	}()
	go func() {
		defer wg.Done()
		for v := range chans[1] {
			time.Sleep(2 * time.Second)
			fmt.Println("Reader 2:", v)
		}
	}()
	wg.Wait()

}
