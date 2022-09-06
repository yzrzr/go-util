# go-util

This is a [Go] language util，contains generic collection interfaces and implementations, and other utilities for working with collections. [quick start][].

## Installation

With [Go module][] support (Go 1.17+)

## Example
list:
```go
package main

import (
	"fmt"

	"github.com/yzrzr/go-util/collect"
)

func main() {
	list := collect.NewArrayList[int](8)
	list.Add(1)
	list.Add(2)
	list.Add(3)
	fmt.Println(list.ToArray()) // [1, 2, 3]
	list.Sort(func(a, b int) bool {
		return a > b
	})
	fmt.Println(list.ToArray()) // [3, 2, 1]
}
```
set:
```go
package main

import (
	"fmt"

	"github.com/yzrzr/go-util/collect"
)

func main() {
	set := collect.NewHashSet[int]()
	set.Add(10)
	set.Add(10)
	set.Add(20)
	fmt.Println(set.ToArray()) // [10, 20]
}

```

## Learn more

- [Examples](examples)

