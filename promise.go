package goncu

import "sync/atomic"

// Promise allows to carry out heavy tasks in parallel
// and get the results asynchronously when needed
type promise[T any] struct {
	responseCh chan T
	errorCh    chan error
	ok         *atomic.Bool
}

type Task[T any] func() (T, error)

// Creates a new promise
func NewPromise[T any](task Task[T]) promise[T] {
	respCh := make(chan T, 1)
	errCh := make(chan error, 1)
	ok := &atomic.Bool{}

	go func() {
		resp, err := task()

		ok.Store(err == nil)
		errCh <- err
		respCh <- resp

		defer close(errCh)
		defer close(respCh)
	}()

	return promise[T]{
		responseCh: respCh,
		errorCh:    errCh,
		ok:         ok,
	}
}

// Set the success case handler
func (p promise[T]) Then(handler func(T)) promise[T] {
	resp := <-p.responseCh

	if p.taskCompleted() {
		handler(resp)
	}

	return p
}

func (p promise[T]) taskCompleted() bool {
	return p.ok.Load()
}

// Set the error case handler
func (p promise[T]) Catch(handler func(error)) promise[T] {
	err := <-p.errorCh

	if err != nil {
		handler(err)
	}

	return p
}

// Take the task result and return it.
// if the task is not ready yet, it will wait until it is ready
func (p promise[T]) Done() (T, error) {
	resp := <-p.responseCh
	err := <-p.errorCh

	return resp, err
}
