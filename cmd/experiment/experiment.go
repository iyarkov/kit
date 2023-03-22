package main

import (
	"fmt"
	"os"
)

func main() {
	for i, a := range os.Args {
		fmt.Printf("argument %d : %s", i, a)
	}
}
