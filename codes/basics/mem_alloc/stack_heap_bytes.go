package main

import "encoding/binary"

func f(x uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, x)
	return buf
}

func g(x uint64) [8]byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], x)
	return buf
}

func main() {
	f(10)
	g(10)
}

/*
go build -gcflags="-m" .
line 6: escapes to heap
*/
