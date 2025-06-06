package basics

import (
	"fmt"
	"maps"
	"testing"
)

type Value struct {
	val int
}

func TestMapOveride(t *testing.T) {
	m := map[int]Value{}
	v := m[0]
	println(v.val)
	v.val = 1
	// cannot override map value directly
	println(v.val, m[0].val)

	mp := map[int]*Value{}
	mp[0] = &Value{val: 0}
	vmp := mp[0]
	vmp.val = 1
	println(vmp.val, mp[0].val)
}

func TestStructWithMapInitalization(t *testing.T) {
	type S struct {
		M map[int]int
	}

	// map 不会自动初始化，
	// s := S{}
	s := S{map[int]int{}}
	s.M[0] = 1
	println(s.M[0])
}

func TestNestedMap(t *testing.T) {
	// 双层map可以自动初始化
	m := map[int]map[int]int{}
	a := m[0]
	fmt.Println("layer1 value:", a)
	b := a[0]
	fmt.Println("layer2 value:", b)

	//
	if a == nil {
		a = map[int]int{}
		m[0] = a
	}
	a[0] = 1
	a[1] = 1
	fmt.Println("m:", m, "a:", a)
}

var All = map[int]map[int]int{1: {}}

type A struct {
	val map[int]int
}

func setMap(a *A) {
	b := All[1]
	b[0] = 1
	a.val = b
}

func TestMapScop(t *testing.T) {
	var a = A{}
	setMap(&a)
	fmt.Println(a)
	fmt.Println(All)

	b := A{}
	b.val = a.val
	fmt.Println(b)
}

// Map ShadowCopy will not change the cloned map when new kv is inserted.
func TestMapClone(t *testing.T) {
	a := map[int]int{1: 2}
	b := maps.Clone(a)
	b[2] = 3
	fmt.Println(a, b)
}
