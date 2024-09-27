package interfaces

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Elem struct {
	Result interface{}
}

func getSingle(out interface{}) {
	val := reflect.ValueOf(out)
	print()
	output := []byte{22, 11}
	val.Elem().Set(reflect.ValueOf(output))
}

func TestReflect(t *testing.T) {
	batch := make([]Elem, 1)
	first := batch[0]
	first.Result = new([]byte)
	fmt.Println("result:", first.Result)
	getSingle(first.Result)
	fmt.Println("result:", first.Result)
	fmt.Println("result val: %", batch[0])
}

func getSingleHexBytes(out interface{}) {
	val := reflect.ValueOf(out)
	output := hexutil.Bytes([]byte{0x11, 0x12})
	val.Elem().Set(reflect.ValueOf(&output))
}

func TestReflectHexBytes(t *testing.T) {
	batch := make([]Elem, 1)
	first := batch[0]
	first.Result = new(hexutil.Bytes)
	fmt.Println("result:", first.Result)
	getSingleHexBytes(&first.Result)
	fmt.Println("result:", first.Result)
	fmt.Println("result val: %", batch[0])
}
