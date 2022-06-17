package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

var keepDaily int

func init() {
	flag.IntVar(&keepDaily, "keep-daily", 0, "daily files/directories to keep")
}

func main() {
	flag.Parse()
	fmt.Println("keep-daily: ", keepDaily)
}
