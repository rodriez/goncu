package goncu

import (
	"errors"
	"time"
)

type poolItem struct {
	data interface{}
	err  error
}

type Pool struct {
	workers     int
	exec        func(n int) (interface{}, error)
	sendSuccess func(item interface{})
	sendError   func(e error)
}

func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}

	return &Pool{
		workers: workers,
	}
}

//Define the function wich is go to be executed
func (p *Pool) DO(fn func(n int) (interface{}, error)) *Pool {
	p.exec = fn
	return p
}

//Define the success listener
func (p *Pool) OnSuccess(fn func(item interface{})) *Pool {
	p.sendSuccess = fn
	return p
}

//Define the error listener
func (p *Pool) OnError(fn func(e error)) *Pool {
	p.sendError = fn
	return p
}

//Run pool execution
func (p *Pool) Run(tasks int) (*PoolResponse, error) {
	if err := p.validate(tasks); err != nil {
		return nil, err
	}

	start := time.Now()
	jobs := make(chan int, tasks)
	results := make(chan poolItem)
	defer close(results)

	p.ready(jobs, results)
	p.start(tasks, jobs)

	response := PoolResponse{}
	for i := 0; i < tasks; i++ {
		item := <-results
		if item.err != nil {
			response.Errors = append(response.Errors, item.err)

			if p.sendError != nil {
				p.sendError(item.err)
			}

			continue
		}

		response.Hits = append(response.Hits, item.data)
		if p.sendSuccess != nil {
			p.sendSuccess(item.data)
		}
	}

	end := time.Now()
	response.Duration = end.Sub(start)

	return &response, nil
}

func (p *Pool) validate(tasks int) error {
	if tasks <= 0 {
		return errors.New("invalid task amount")
	}

	if p.exec == nil {
		return errors.New("there is nothing to do")
	}

	return nil
}

func (p *Pool) ready(jobs <-chan int, results chan<- poolItem) {
	for id := 0; id < p.workers; id++ {
		go p.worker(id, jobs, results)
	}
}

func (p *Pool) worker(wokerId int, jobs <-chan int, results chan<- poolItem) {
	for n := range jobs {
		data, e := p.exec(n)
		results <- poolItem{
			data: data,
			err:  e,
		}
	}
}

func (p *Pool) start(tasks int, jobs chan<- int) {
	for i := 0; i < tasks; i++ {
		jobs <- i
	}
	close(jobs)
}
