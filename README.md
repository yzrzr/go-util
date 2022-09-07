# go-util

This is a [Go] language util，contains generic collection interfaces and implementations, and other utilities for working with collections. [quick start][].

## Installation

With [Go module][] support (Go 1.17+)

## Interface

- [Collection](collect/collection.go)
- [List](collect/list.go)
- [Set](collect/set.go)
- [Iterator](collect/iterator.go)

## Example
list:
```go
package main

import (
	"fmt"

	"github.com/yzrzr/go-util/collect"
)

func main() {
	list := collect.NewList[int](collect.DefaultListConfig)
	list.Add(1)
	list.Add(2)
	list.Add(3)
	fmt.Println(list.ToArray()) // [1, 2, 3]
	list.Sort(collect.SortLessOrdered[int](false))
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
	set := collect.NewSet[int]()
	set.Add(10)
	set.Add(10)
	set.Add(20)
	fmt.Println(set.ToArray()) // [10, 20]
}

```

## Learn more

- [Examples](examples)

