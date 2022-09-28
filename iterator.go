package goncu

type Producer[T any] func(i int) T

// Iterator is used to make a function iterable
//
// This helps optimize long processing times, since while Producer is busy,
// an external consumer can process the iterations each returns, in parallel.
type iterator[T any] struct {
	// indicates the number of times Producer will run
	length int
	// Function that receives an iteration number and returns an element to iterate through Each
	produce Producer[T]
}

// Create a new iterator.
// length - indicates the number of times Producer will run
// producer - Function that receives an iteration number and returns an element to iterate through Each
func Iterator[T any](length int, producer Producer[T]) iterator[T] {
	return iterator[T]{
		length:  length,
		produce: producer,
	}
}

// Create a new iterator from a slice.
func SliceIterator[T any](s []T) iterator[T] {
	return iterator[T]{
		length: len(s),
		produce: func(i int) T {
			return s[i]
		},
	}
}

// Returns a channel that reads one element per iteration.
// Should be used as a For-range loop target
func (g iterator[T]) Each() chan T {
	c := make(chan T)

	if g.produce == nil {
		close(c)
		return c
	}

	go func() {
		for i := 0; i < g.length; i++ {
			c <- g.produce(i)
		}

		close(c)
	}()

	return c
}

// Execute a function for every element that the iterator produce.
func (g iterator[T]) ForEach(handler func(int, T)) {
	i := 0
	for v := range g.Each() {
		handler(i, v)
		i++
	}
}

// Reduce the iterator to an array of elements.
func (g iterator[T]) ToArray() []T {
	array := make([]T, g.length)

	g.ForEach(func(i int, s T) {
		array[i] = s
	})

	return array
}
