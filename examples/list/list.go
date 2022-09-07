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
