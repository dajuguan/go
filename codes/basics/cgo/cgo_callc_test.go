package main

import (
	"testing"
)

func TestCallCFromGo(t *testing.T) {
	Random()
	Add(1, 2)
}
