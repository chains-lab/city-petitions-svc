package main

import (
	"os"

	"github.com/chains-lab/city-petitions-svc/cmd/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
