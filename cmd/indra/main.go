package main

import (
	"errors"
	"os"

	indra "coderaiser/indra"
)

func main() {
	if err := indra.Indra(Registry, os.Args[1:], os.Stderr); err != nil {
		if errors.Is(err, indra.ErrInvalidOption) {
			os.Exit(7)
		}
		os.Exit(1)
	}
}
