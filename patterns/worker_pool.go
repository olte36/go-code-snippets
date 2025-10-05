package patterns

import (
	"errors"
	"log"
	"sync/atomic"
)

var ErrPoolIsClosed = errors.New("the pool is closed")

type Result[T any] struct {
	Res any
	Err error
}

type WorkerPool[T any] struct {
	jobs    chan func() (T, error)
	results chan Result[T]
	closed  atomic.Bool
}

func NewWorkerPool[T any](numWorkes int) *WorkerPool[T] {
	wp := WorkerPool[T]{
		jobs:    make(chan func() (T, error)),
		results: make(chan Result[T]),
	}
	for i := range numWorkes {
		go wp.work(i)
	}
	return &wp
}

func (w *WorkerPool[T]) Go(f func() (T, error)) error {
	if w.closed.Load() {
		return ErrPoolIsClosed
	}
	w.jobs <- f
	return nil
}

func (w *WorkerPool[T]) TryGo(f func() (T, error)) (bool, error) {
	if w.closed.Load() {
		return false, ErrPoolIsClosed
	}
	select {
	case w.jobs <- f:
		return true, nil
	default:
		return false, nil
	}
}

func (w *WorkerPool[T]) Results() <-chan Result[T] {
	return w.results
}

func (w *WorkerPool[T]) Close() {
	w.closed.Store(true)
	close(w.jobs)
	close(w.results)
}

func (w *WorkerPool[T]) work(id int) {
	log.Printf("worker #%d started", id)
	for job := range w.jobs {
		res, err := job()
		w.results <- Result[T]{
			Res: res,
			Err: err,
		}
		log.Printf("worker #%d finished job", id)
	}
	log.Printf("worker #%d stopped", id)
}
