# goncu
A simple library of conccurrent tools

## Installation

```bash
go get github.com/rodriez/goncu
```

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

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

[Hey 👋 buy me a beer! ](https://www.buymeacoffee.com/rodriez)