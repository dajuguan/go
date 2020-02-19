package store

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

//Node 节点
type Node struct {
	Value  []byte
	Childs map[rune]*Node
}

//Trie  树
type Trie struct {
	RootHash []byte
	Head     *Node
}

// type Trie interface {
// 	Get(key []byte) []byte。h
// 	Update(key, value []byte)
// 	Delete(key []byte)
// 	Iter()
// 	GetStateRoot() []byte
// }

//SerializeTrie 序列化
func (t *Trie) SerializeTrie() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(t)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

//DeserializeTrie 反序列化
func DeserializeTrie(data []byte) *Trie {
	var t Trie
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	return &t
}

func (t *Trie) GenRootHash() {
	hash := sha256.Sum256(t.SerializeTrie())
	t.RootHash = hash[:]
}

func (t *Trie) Update(key string, value []byte) {
	node := t.Head
	for _, char := range key {
		if node.Childs[char] == nil {
			node.Childs[char] = &Node{[]byte{}, make(map[rune]*Node)}
		}
		node = node.Childs[char]
	}
	node.Value = value
	t.GenRootHash()
}

func (t *Trie) Get(key string) (*Node, error) {
	node := t.Head
	for _, char := range key {
		if node.Childs[char] != nil {
			node = node.Childs[char]
		} else {
			return nil, fmt.Errorf("获取的字符%s不存在", key)
		}
	}
	return node, nil
}

func InitTrie() *Trie {
	head := &Node{[]byte{}, make(map[rune]*Node)}
	t := Trie{
		Head:     head,
		RootHash: []byte("111111"),
	}
	return &t
}
