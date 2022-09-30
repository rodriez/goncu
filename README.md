# goncu
A simple library with concurrent tools

[![Go Report Card](https://goreportcard.com/badge/github.com/rodriez/goncu/v2)](https://goreportcard.com/report/github.com/rodriez/goncu/v2)   [![PkgGoDev](https://pkg.go.dev/badge/github.com/rodriez/goncu/goncu/v2)](https://pkg.go.dev/github.com/rodriez/goncu/v2)

## Installation

```bash
go get github.com/rodriez/goncu/v2
```

## Iterator Usage
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu/v2"
)

func main() {
	iterator := goncu.Iterator(100, func(i int) string {
		return fmt.Sprintf("%d", i)
	})

	newSlice := []string{}
	for e := range iterator.Each() {
		newSlice = append(newSlice, e)
	}

	fmt.Println(newSlice)
}
```

## Iterator ForEach 
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu/v2"
)

func main() {
	iterator := goncu.Iterator(100, func(i int) string {
		return fmt.Sprintf("%d", i)
	})

	newSlice := []string{}
	iterator.ForEach(func(i int, s string) {
		newSlice = append(newSlice, s)
	})

	fmt.Println(newSlice)
}
```

## Iterator To Array 
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu/v2"
)

func main() {
	iterator := goncu.Iterator(100, func(i int) string {
		return fmt.Sprintf("%d", i)
	})

	newSlice := iterator.ToArray()

	fmt.Println(newSlice)
}
```

## Slice Iterator Usage
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu/v2"
)

func main() {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	iterator := goncu.SliceIterator(data)

	newSlice := []string{}
	for e := range iterator.Each() {
		newSlice = append(newSlice, e)
	}

	fmt.Println(newSlice)
}
```

## Basic Worker Pool Usage
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu/v2"
)

func main() {
	hits := int32(0)
	err := goncu.Pool(2, 10).Run(func(n int) {
		atomic.AddInt32(&hits, n)
	})

    fmt.Println(hits)
    fmt.Println(err)
}
```

## Iterator Worker Pool Usage
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu/v2"
)

func main() {
	iterator := goncu.Iterator(30, func(i int) string {
		return fmt.Sprintf("%d", i)
	})

	buffer := make(chan string, 100)
	defer close(buffer)
	err := goncu.IteratorPool(10, iterator).Run(func(s string) {
		buffer <- s
	})

	for hit := range buffer {
		fmt.Print(hit)
	}
	fmt.Print("\n")
    fmt.Println(err)
}
```

## Slice Worker Pool Usage
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu/v2"
)

func main() {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	buffer := make(chan string, 100)
	defer close(buffer)
	err := goncu.SlicePool(-1, data).Run(func(s string) {
		buffer <- s
	})

	for hit := range buffer {
		fmt.Print(hit)
	}
	fmt.Print("\n")
    fmt.Println(err)
}
```

## Sync Promise Usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/rodriez/goncu/v2"
)

func main() {
	promise := goncu.NewPromise(func() (string, error) {
		time.Sleep(50 * time.Millisecond)
		return "A ver long task was executed", nil
	})

	fmt.Println("Do a lot of things")
	time.Sleep(100 * time.Millisecond)

	str, err := promise.Done()
	
   	fmt.Println(str)
   	fmt.Println(err)
}
```

## Async Promise Usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/rodriez/goncu/v2"
)

func main() {
	goncu.NewPromise(func() (string, error) {
		time.Sleep(50 * time.Millisecond)

		return "A ver long task was executed", nil
	}).Then(func(str string) {
		fmt.Println(str)
	}).Catch(func(err error) {
		fmt.Println(err)
	})

	time.Sleep(100 * time.Millisecond)
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

[Hey 👋 buy me a beer! ](https://www.buymeacoffee.com/rodriez)
