package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const name = "Hello, OTUS!"

func main() {
	fmt.Printf("%s", stringutil.Reverse(name))
}
