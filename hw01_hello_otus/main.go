package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	str := "Hello, OTUS!"
	reversed := reverse.String(str)
	fmt.Println(reversed)
}
