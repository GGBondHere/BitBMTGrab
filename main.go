package main

import (
	"badminton/app"
	"fmt"
	"os"
)

var version = "dev"

func main() {
	err := app.Run(version)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
