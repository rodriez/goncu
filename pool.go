// Package goncu provides structures that simplify the implementation
// of solutions for concurrent problems
package goncu

import (
	"fmt"
	"sync"
)

// Action to execute
type Job[T any] func(T)

// Pool allows the distribution of a set of tasks among a certain number of workers
type pool[T any] struct {
	size     int
	iterator iterator[T]
}

// Creates a new workers Pool from an iterator
//
// size - indicates the number of workers that the Pool will use to manage the set of tasks,
// the minimum value allowed is 1 and the maximum is the same number of tasks
// iterator - iterator thar produce the items to process
func IteratorPool[T any](size int, iterator iterator[T]) pool[T] {
	if iterator.length < size || size <= 0 {
		size = 1
	}

	return pool[T]{
		size:     size,
		iterator: iterator,
	}
}

// Creates a new workers Pool from a slice
//
// size - indicates the number of workers that the Pool will use to manage the set of tasks,
// the minimum value allowed is 1 and the maximum is the same number of tasks
// s - slice with items to process
func SlicePool[T any](size int, s []T) pool[T] {
	iterator := SliceIterator(s)

	if iterator.length < size || size <= 0 {
		size = 1
	}

	return pool[T]{
		size:     size,
		iterator: iterator,
	}
}

// Creates a new workers Pool
//
// size - indicates the number of workers that the Pool will use to manage the set of tasks,
// the minimum value allowed is 1 and the maximum is the same number of tasks
// iterations - max number of task to process
func Pool(size int, iterations int) pool[int] {
	iterator := Iterator(iterations, func(i int) int {
		return i
	})

	if iterator.length < size || size <= 0 {
		size = 1
	}

	return pool[int]{
		size:     size,
		iterator: iterator,
	}
}

// Run start pool execution, collect all data and return a response
// job - Action to execute
func (p pool[T]) Run(job Job[T]) error {
	if job == nil {
		return fmt.Errorf("invalid pool job")
	}

	if p.iterator.length <= 0 {
		return fmt.Errorf("invalid number of iterations")
	}

	var wg sync.WaitGroup
	wg.Add(p.iterator.length)

	jobs := make(chan T)
	defer close(jobs)

	p.ready(jobs, job, &wg)
	p.start(jobs)

	wg.Wait()

	return nil
}

func (p pool[T]) ready(jobs <-chan T, job Job[T], wg *sync.WaitGroup) {
	for i := 0; i < p.size; i++ {
		go p.launchWorker(jobs, job, wg)
	}
}

func (p pool[T]) launchWorker(jobs <-chan T, job Job[T], wg *sync.WaitGroup) {
	for j := range jobs {
		job(j)
		wg.Done()
	}
}

func (p pool[T]) start(jobs chan<- T) {
	for job := range p.iterator.Each() {
		jobs <- job
	}
}
