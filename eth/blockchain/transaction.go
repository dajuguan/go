package blockchain

import (
	"fmt"
)

type TxType int8

const (
	TxNormal         TxType = 0
	TxCoinbase       TxType = 1
	TxContractCreate TxType = 2
	TxContractCall   TxType = 3
)

type Transaction struct {
	Id       []byte
	From     []byte
	To       []byte
	Value    int
	Data     []byte
	Sign     []byte
	gasLimit int
	gasPrice int
	Type     TxType
	nonce    int
}

func CreateTx(from, to []byte, value int) *Transaction {
	return &Transaction{
		From:  from,
		To:    to,
		Value: value,
		Type:  TxNormal,
	}
}

func CoinbaseTx(to []byte, value int) *Transaction {
	data := fmt.Sprintf("%x", "Coinbase Tx")
	return &Transaction{
		To:    to,
		Value: value,
		Data:  []byte(data),
		Type:  TxCoinbase,
	}
}

func ContractCreateTx() *Transaction {
	return &Transaction{
		Data: []byte{},
		Type: TxContractCreate,
	}
}
