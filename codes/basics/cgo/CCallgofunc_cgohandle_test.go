package main

import (
	"testing"
)

func TestCCallGoWithCgoHandle(t *testing.T) {
	ExampleCallHandle()
	ExampleCallStructHandle()
}
