package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/dajuguan/go/codes/wallet"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// Coinbase交易
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 24)
		_, err := rand.Read(randData)
		Handle(err)
		data = fmt.Sprintf("%x", randData)
	}

	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(100, to)
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.SetID()

	return &tx
}

// 设置ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Outputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// map[hash]Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//coninbase不需要签名
	if tx.IsCoinbase() {
		return
	}
	//检查输入是否合法
	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Error:引用的交易不存在")
		}
	}
	txCopy := tx.TrimmedCopy()
	for inId, in := range txCopy.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].Pubkey = prevTx.Outputs[in.Out].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].Pubkey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Inputs[inId].Signature = signature
	}
}

// 实现深拷贝，防止变量被修改
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}
	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubkeyHash})
	}
	txCopy := Transaction{ID: tx.ID, Inputs: inputs, Outputs: outputs}
	return txCopy
}

func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, in := range tx.Inputs {
		if prevTxs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Error:引用的交易不存在")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for inID, in := range tx.Inputs {
		prevTx := prevTxs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].Pubkey = prevTx.Outputs[in.Out].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inID].Pubkey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:sigLen/2])
		s.SetBytes(in.Signature[sigLen/2:])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.Pubkey)
		x.SetBytes(in.Pubkey[:keyLen/2])
		y.SetBytes(in.Pubkey[(keyLen / 2):])

		rawPubkey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubkey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true

}

func NewTransaction(w *wallet.Wallet, to string, amount int, UTXO *UTXOSet) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	//接入钱包

	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	total, validOutputs := UTXO.FindSpendableOutputs(pubKeyHash, amount)
	if total < amount {
		log.Panic("Error Not Enough Found")
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)
		for _, out := range outs {
			input := TxInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTxOutput(amount, to))

	if total > amount {
		outputs = append(outputs, *NewTxOutput(total-amount, string(w.Address())))
	}

	tx := Transaction{nil, inputs, outputs}
	//传入Hash
	tx.ID = tx.Hash()
	//签名
	UTXO.BC.SignTransaction(&tx, w.PrivateKey)

	return &tx
}

func DeserializeTransaction(data []byte) Transaction {
	var tx Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(tx)
	if err != nil {
		Handle(err)
	}
	return tx
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 为什么要采用复制的方式？
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("-- Transaction %x:", tx.ID))
	for i, in := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("    Input %d:", i))
		lines = append(lines, fmt.Sprintf("    TxID:      %x", in.ID))
		lines = append(lines, fmt.Sprintf("    Out:       %d", in.Out))
		lines = append(lines, fmt.Sprintf("    Signature: %x", in.Signature))
		lines = append(lines, fmt.Sprintf("    PubKey:    %x", in.Pubkey))
	}
	for i, out := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("    Output %d:", i))
		lines = append(lines, fmt.Sprintf("    Value:      %x", out.Value))
		lines = append(lines, fmt.Sprintf("    PubKey:       %d", out.PubkeyHash))
	}
	return strings.Join(lines, "\n")
}
