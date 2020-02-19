//账户功能
//1.创建账户
//2.签名交易
//3.存储账户状态

package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
)

type Wallet struct {
	Accounts map[string]*Account
}

func InitWallet() (*Wallet, error) {
	wallet := Wallet{}
	wallet.Accounts = make(map[string]*Account)
	return &wallet, nil
}

func (w *Wallet) AddAccount() string {
	account := NewAccount()
	addr := fmt.Sprintf("添加账户地址:%s", account.Address())
	w.Accounts[addr] = account
	return addr
}

func (w *Wallet) GetAddr(addr string) *Account {
	return w.Accounts[addr]
}

func (w *Wallet) GetAllAddr() []string {
	var addrs []string
	for addr := range w.Accounts {
		addrs = append(addrs, addr)
	}
	return addrs
}

func SignTx(account *Account, txID []byte) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, &account.PrivateKey, txID)
	if err != nil {
		log.Panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}
