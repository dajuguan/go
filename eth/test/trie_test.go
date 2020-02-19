package test

import (
	"bytes"
	"testing"

	"github.com/dajuguan/go/eth/store"
)

func TestTrie(t *testing.T) {
	trie := store.InitTrie()
	t.Log(trie.RootHash)
	trie.Update("a", []byte("first"))
	t.Log(trie.RootHash)
	node, _ := trie.Get("a")
	if bytes.Compare(node.Value, []byte("first")) != 0 {
		t.Errorf("trie.Get(a) 应该为first")
	}
	node, _ = trie.Get("b")
	if node != nil {
		t.Errorf("trie.Get(b) 应该为nil")
	}
	trie.Update("ab", []byte("second"))
	t.Log(trie.RootHash)
	node, _ = trie.Get("ab")
	if bytes.Compare(node.Value, []byte("second")) != 0 {
		t.Errorf("trie.Get(a) 应该为second")
	}
}
