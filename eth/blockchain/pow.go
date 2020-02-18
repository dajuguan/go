package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

type Pow struct {
	Block  *Block
	Target *big.Int
}

func NewPow(block *Block) *Pow {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-block.Difficulty))
	return &Pow{block, target}
}

func ToHex(x int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, x)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()

}

func (pow *Pow) InitData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.ParentHash,
		pow.Block.ReceiptRoot,
		pow.Block.StateRoot,
		pow.Block.TxRoot,
		ToHex(pow.Block.Timestamp),
		ToHex(int64(nonce)),
		ToHex(int64(pow.Block.Difficulty)),
	}, []byte{})
	return data[:]
}

func (pow *Pow) run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hashed := sha256.Sum256(data)
		hashInt.SetBytes(hashed[:])
		fmt.Printf("\r%x", hashed)
		if hashInt.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}
	fmt.Println()
	return nonce, hash[:]
}

//验证是否满足难度条件
func ValidatePow(b *Block) bool {
	pow := NewPow(b)
	data := pow.InitData(b.Nonce)
	hash := sha256.Sum256(data)
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.Target) == -1
}
