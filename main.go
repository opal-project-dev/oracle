package main

import (
	"os"

	"github.com/opal-project-dev/oracle/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
