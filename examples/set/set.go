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
