package main

import (
	"os"

	"bragg-xrd/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
