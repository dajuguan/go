package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/dajuguan/go/codes/wallet"
)

type TxOutput struct {
	Value      int
	PubkeyHash []byte
}

type TxInput struct {
	ID        []byte //关键的交易
	Out       int    //关联输出的index
	Signature []byte
	Pubkey    []byte
}

type TxOutputs struct {
	Outputs []TxOutput
}

func NewTxOutput(value int, address string) *TxOutput {
	txo := TxOutput{value, nil}
	txo.Lock([]byte(address))
	return &txo
}

func (in *TxInput) UserKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.Pubkey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubkeyHash = pubKeyHash
}

func (out *TxOutput) isLockedWithPubKey(pubKeyHash []byte) bool {
	return bytes.Compare(pubKeyHash, out.PubkeyHash) == 0
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(outs)
	Handle(err)
	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outs TxOutputs
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outs)
	Handle(err)
	return outs
}
