package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	hello := "Hello, OTUS!"
	hello = stringutil.Reverse(hello)
	fmt.Println(hello)
}
