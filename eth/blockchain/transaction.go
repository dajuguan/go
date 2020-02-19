package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"math/big"
)

type TxType int8

const (
	TxNormal         TxType = 0
	TxCoinbase       TxType = 1
	TxContractCreate TxType = 2
	TxContractCall   TxType = 3
)

type Transaction struct {
	Id        []byte
	From      []byte
	To        []byte
	Value     int
	Data      []byte
	Sign      []byte
	gasLimit  int
	gasPrice  int
	Type      TxType
	nonce     int
	Signature []byte
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

//验证交易输入的签名
func VerifySig(signature, pubKeyHash, data []byte) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:sigLen/2])
	s.SetBytes(signature[sigLen/2:])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubKeyHash)
	x.SetBytes(pubKeyHash[:keyLen/2])
	y.SetBytes(pubKeyHash[(keyLen / 2):])
	curve := elliptic.P256()
	rawPubkey := ecdsa.PublicKey{curve, &x, &y}
	return ecdsa.Verify(&rawPubkey, data, &r, &s)
}

func DeserializeTransaction(data []byte) Transaction {
	var tx Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(tx)
	if err != nil {
		panic(err)
	}
	return tx
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		panic(err)
	}
	return encoded.Bytes()
}
