package goncu

// Promise allows to carry out heavy tasks in parallel
// and get the results asynchronously when needed
type promise[T any] struct {
	ch  chan T
	err chan error
}

type Task[T any] func() (T, error)

// Creates a new promise
func NewPromise[T any](task Task[T]) promise[T] {
	p := promise[T]{
		ch:  make(chan T, 1),
		err: make(chan error, 1),
	}

	go func() {
		resp, err := task()

		if err != nil {
			p.err <- err
		} else {
			p.ch <- resp
		}

		defer p.close()
	}()

	return p
}

func (p promise[T]) Then(handler func(T)) promise[T] {
	resp := <-p.ch
	handler(resp)

	return p
}

func (p promise[T]) Catch(handler func(error)) promise[T] {
	err := <-p.err
	handler(err)

	return p
}

func (p promise[T]) close() {
	close(p.err)
	close(p.ch)
}

// Take the task result and return it.
// if the task is not ready yet, it will wait until it is ready
func (p promise[T]) Done() (T, error) {
	resp := <-p.ch
	err := <-p.err

	return resp, err
}
