package main

import (
	"os"

	indra "coderaiser/indra"
)

func main() {
	if err := indra.Indra(os.Args[1:], os.Stderr); err != nil {
		os.Exit(1)
	}
}
