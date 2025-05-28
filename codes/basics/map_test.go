package basics

import "testing"

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
