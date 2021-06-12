package goncu

//Iterator is used to make a function iterable
//
//This helps optimize long processing times, since while Producer is busy,
//an external consumer can process the iterations each returns, in parallel.
type Iterator struct {
	length int
	get    func(i int) interface{}
}

//Create a new iterator.
//
//	length - indicates the number of times Producer will run
func NewIterator(length int) *Iterator {
	return &Iterator{
		length: length,
	}
}

//Producer initializes the function to be executed in each iteration.
//
//fn - Function that receives an iteration number and returns an element to iterate through Each
func (g *Iterator) Producer(fn func(i int) interface{}) *Iterator {
	g.get = fn
	return g
}

//Each returns a channel that reads one element per iteration.
//Should be used as a For-range loop target
func (g *Iterator) Each() chan interface{} {
	c := make(chan interface{})

	if g.get == nil {
		close(c)
		return c
	}

	go func() {
		for i := 0; i < g.length; i++ {
			c <- g.get(i)
		}
		close(c)
	}()

	return c
}
