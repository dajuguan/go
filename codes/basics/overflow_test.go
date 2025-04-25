package basics

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
)

// Test for overflow
func TestOverflow(t *testing.T) {
	const hashBytes = 2
	// if size exceeds 4Gb(2^32), index might overflow if limit is u32
	size := uint64(math.MaxUint32 * hashBytes)
	dataSet := make([]byte, size) // 8Gb
	// error limit and overflowed index
	limit := uint32(size / hashBytes)
	indexStart := (limit - 2) * hashBytes

	// expected limit and index
	expectedLimit := uint64(size / hashBytes)
	expectedIndexStart := (expectedLimit - 2) * hashBytes

	dataSet[indexStart] = 1
	dataSet[expectedIndexStart] = 2
	fmt.Println("MaxUint32       :", math.MaxUint32)
	fmt.Println("overflowed index:", indexStart, "val:", dataSet[indexStart])
	fmt.Println("expected index  :", expectedIndexStart, "expectedVal:", dataSet[expectedIndexStart])
}
