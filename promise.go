package goncu

//Promise allows to carry out heavy tasks in parallel
// and get the results asynchronously when needed
type Promise struct {
	ch chan promiseData
}

//Creates a new promise
func NewPromise() *Promise {
	return &Promise{
		ch: make(chan promiseData, 1),
	}
}

//Start passed task that in parallel
//
//task - is a funcion that make a process and return a response or an error
func (p *Promise) Start(task func() (interface{}, error)) *Promise {
	go func() {
		d, e := task()
		p.ch <- promiseData{Data: d, Error: e}
	}()

	return p
}

//Take the task result and return it.
//if the task is not ready yet, it will wait until it is ready
func (p *Promise) Done() (interface{}, error) {
	resp := <-p.ch
	close(p.ch)

	return resp.Data, resp.Error
}

type promiseData struct {
	Data  interface{}
	Error error
}
