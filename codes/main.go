package main

import (
	"os"

	"github.com/dajuguan/go/cli"
)

type test struct {
	name string
}

func nilPtr(a *test) {

}

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
	// var a *test
	// fmt.Println(a, a.name)
}
