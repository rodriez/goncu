# goncu
A simple library with concurrent tools

## Installation

```bash
go get github.com/rodriez/goncu
```

## Worker Pool Usage
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu"
)

func main() {
	response, err := goncu.NewPool(2).DO(func(n int) (interface{}, error) {
		return n + 1, nil
	}).OnSuccess(func(item interface{}) {
		fmt.Println(item)
	}).OnError(func(e error) {
		fmt.Println(e)
	}).Run(10)

    fmt.Println(response)
    fmt.Println(err)
}
```

## Promise Usage
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu"
)

func main() {
	promise := goncu.NewPromise().
		Start(func() (interface{}, error) {
			return "A ver long string", nil
		})

	time.Sleep(100 * time.Millisecond)

	resp, err := promise.Done()

    fmt.Println(response)
    fmt.Println(err)
}
```

## Iterator Usage
```go
package main

import (
	"fmt"

	"github.com/rodriez/goncu"
)

func main() {
	iterator := goncu.NewIterator(100).
		Producer(func(i int) interface{} {
			return fmt.Sprintf("%d", i)
		})

	newSlice := []string{}
	for e := range iterator.Each() {
		newSlice = append(newSlice, e.(string))
	}

    fmt.Println(newSlice)
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

[Hey 👋 buy me a beer! ](https://www.buymeacoffee.com/rodriez)
