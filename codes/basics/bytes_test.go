package basics

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestBytes(t *testing.T) {
	var a []byte
	b := common.Hash{}
	c := common.BytesToHash(a)
	fmt.Println(c == b)
}
