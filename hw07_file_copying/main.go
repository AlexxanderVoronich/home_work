package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" || to == "" {
		fmt.Println("Usage: go run main.go -from <source_file> -to <destination_file> [-offset <offset>] [-limit <limit>]")
		return
	}
	fmt.Printf("Launch with parameters: -from %s -to %s -offset %d -limit %d\n", from, to, offset, limit)

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
