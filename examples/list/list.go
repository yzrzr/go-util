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
