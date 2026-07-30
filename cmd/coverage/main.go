package main

import (
	"errors"
	"fmt"
	"os"

	"coderaiser/indra"
)

func main() {
	if err := indra.Run(os.Args[1:], os.Stdout); err != nil {
		if errors.Is(err, indra.ErrUncovered) {
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
