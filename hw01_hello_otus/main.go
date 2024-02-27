package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	s := "Hello, OTUS!"
	reversed := reverse.String(s)
	fmt.Println(reversed)
}
