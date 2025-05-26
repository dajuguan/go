package basics

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Person struct {
	Name *string `json:"name"`
	Age  *int    `json:"age"`
	Root [2]byte
}

func TestJsonPointer(t *testing.T) {
	jsonStr := `{"name": "Alice", "age": null}`

	var p Person
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		panic(err)
	}
	pp := p
	p.Root[0] = 1

	fmt.Println("Name:", *p.Name) // 非 nil
	fmt.Println("Age:", p.Age)    // nil，因为 JSON 里是 null

	fmt.Println("root:", pp.Root, p.Root)
	// save Person to json
	res, _ := json.Marshal(p)
	fmt.Println(string(res))
}
