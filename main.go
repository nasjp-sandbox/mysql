package main

import (
	"fmt"
	"os"
)

func main() {
	// if err := transactionOrder(); err != nil {
	if err := forUpdate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
