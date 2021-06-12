//Package goncu provides structures that simplify the implementation
//of solutions for concurrent problems
package goncu

import (
	"errors"
	"time"
)

//Pool allows the distribution of a set of tasks among a certain number of workers
type Pool struct {
	workers     int
	task        func(n int) (interface{}, error)
	sendSuccess func(item interface{})
	sendError   func(e error)
}

//Creates a new workers Pool
//
//workers - indicates the number of workers that the Pool will use to manage the set of tasks,
//the minimum value allowed is 1 and the maximum is the same number of tasks
func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}

	return &Pool{
		workers: workers,
	}
}

//Define the function wich is go to be executed
//
//task - Function that receives an iteration number and returns an reponse or an error
func (p *Pool) DO(task func(n int) (interface{}, error)) *Pool {
	p.task = task
	return p
}

//Define the success event listener
//
//event - Function that receives the task response
func (p *Pool) OnSuccess(event func(item interface{})) *Pool {
	p.sendSuccess = event
	return p
}

//Define the error event listener
//
//event - Function that receives the task error
func (p *Pool) OnError(event func(e error)) *Pool {
	p.sendError = event
	return p
}

//Run start pool execution, collect all data and return a response
//
//times - Indicates the number of times than task should be called
func (p *Pool) Run(times int) (*PoolResponse, error) {
	if err := p.validate(times); err != nil {
		return nil, err
	}

	if p.workers > times {
		p.workers = times
	}

	start := time.Now()
	jobs := make(chan int, times)
	results := make(chan poolItem)
	defer close(results)

	p.ready(jobs, results)
	p.start(times, jobs)

	response := p.buildResponse(results, times)
	response.Duration = time.Since(start)
	return response, nil
}

func (p *Pool) validate(times int) error {
	if times <= 0 {
		return errors.New("invalid task amount")
	}

	if p.task == nil {
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
		data, e := p.task(n)

		if e != nil && p.sendError != nil {
			p.sendError(e)
		} else if data != nil && p.sendSuccess != nil {
			p.sendSuccess(data)
		}

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

func (p *Pool) buildResponse(results <-chan poolItem, items int) *PoolResponse {
	response := PoolResponse{}

	for i := 0; i < items; i++ {
		if item := <-results; item.err != nil {
			response.Errors = append(response.Errors, item.err)
		} else {
			response.Hits = append(response.Hits, item.data)
		}
	}

	return &response
}

//PoolResponse Contains the pool process duration also the hits and errors collected
type PoolResponse struct {
	Duration time.Duration
	Errors   []error
	Hits     []interface{}
}

type poolItem struct {
	data interface{}
	err  error
}
