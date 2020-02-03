package main

import (
	"os"

	"github.com/dajuguan/go/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
