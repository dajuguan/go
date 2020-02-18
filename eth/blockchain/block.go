package blockchain

import (
	"bytes"
	"log"
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

func Validate(lastBlock, b *Block, state State) bool {
	//创世区块
	if bytes.Compare(b.ParentHash, []byte{}) == 0 {
		return true
	}

	if bytes.Compare(b.ParentHash, lastBlock.Hash) != 0 {
		log.Panic("Hash与父区块不一致")
	}
	if b.Height != lastBlock.Height+1 {
		log.Panic("区块高度不符合")
	}
	//验证交易根
	//验证难度
	if !ValidatePow(b) {
		log.Panic("挖矿难度不符合")
	}

	return true
}
