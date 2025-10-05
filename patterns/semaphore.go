package patterns

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(maxGoroutines int) *Semaphore {
	if maxGoroutines < 1 {
		panic("max goroutines must be >= 1")
	}
	return &Semaphore{
		ch: make(chan struct{}, maxGoroutines),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}
