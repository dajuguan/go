package main

import (
	"fmt"

	"github.com/dajuguan/go/eth/blockchain"
)

func main() {
	print("test")
	chain := blockchain.InitBlockChain()
	genesis := chain.Blocks[0]
	firstBlock := blockchain.CreateBlock(genesis.Hash, []byte{}, []byte{}, []byte{}, blockchain.DIFFICULTY, 1)
	chain.AddBlock(firstBlock)
	fmt.Println(chain)
}
