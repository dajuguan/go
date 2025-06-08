package basics

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type PointerStruct struct {
	p *int
}

func TestPointerChange(t *testing.T) {
	v := new(int)
	*v = 1
	a := PointerStruct{p: v}
	b := a
	a.p = new(int)
	fmt.Println(*a.p, *b.p)

}

func TestPointerByte(t *testing.T) {
	var v []byte
	fmt.Println(v == nil)
}

func TestNilEqualBytes(t *testing.T) {
	a := common.Hash{}
	var b common.Hash
	fmt.Println(a == b)
}
