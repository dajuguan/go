package blockchain

import (
	"time"
)

const DIFFICULTY = 14

type Block struct {
	Hash        []byte
	ParentHash  []byte
	StateRoot   []byte
	TxRoot      []byte
	ReceiptRoot []byte
	Difficulty  int
	Timestamp   int64
	Height      int
	// Transactions []*Transaction
	Nonce int
}

func CreateBlock(parentHash, stateRoot, txRoot, receiptRoot []byte, diff, height int) *Block {
	block := &Block{
		ParentHash:  parentHash,
		StateRoot:   stateRoot,
		TxRoot:      txRoot,
		ReceiptRoot: receiptRoot,
		Difficulty:  diff,
		Height:      height,
		Timestamp:   time.Now().Unix(),
	}
	pow := NewPow(block)
	nonce, hash := pow.run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func Genesis() *Block {
	return CreateBlock([]byte{}, []byte{}, []byte{}, []byte{}, DIFFICULTY, 0)
}
