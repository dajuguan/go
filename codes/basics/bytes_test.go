package basics

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestBytes(t *testing.T) {
	var a []byte
	b := common.Hash{}
	c := common.BytesToHash(a)
	fmt.Println(c == b)
}

func TestBytesNil(t *testing.T) {
	a := []byte{}
	// false
	fmt.Println(a == nil)

	var b []byte
	// true
	fmt.Println(b == nil)
	hashNil := crypto.Keccak256Hash(a)
	hashNilBytes := crypto.Keccak256Hash(b)
	// true
	fmt.Println(hashNil.Hex(), hashNilBytes.Hex())
}
