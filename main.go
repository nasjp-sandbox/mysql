package main

import (
	"fmt"
	"os"

	"github.com/nasjp-sandbox/mysql/cases"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	if err := cases.LostUpdate(); err != nil {
		return err
	}

	if err := cases.ForUpdate(); err != nil {
		return err
	}

	return nil
}
