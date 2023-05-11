package main

import (
	"fmt"
)

type typeA struct {
	value string
}

func main() {
	sl := map[int32]typeA{
		0: {value: "1"},
	}
	fmt.Printf("Value [0]: %v\n", sl[0])
	fmt.Printf("Value [1]: %v\n", sl[1])
}
