package test

import (
	"testing"

	"github.com/dajuguan/go/eth/wallet"
)

func TestWallet(t *testing.T) {
	wallet, _ := wallet.InitWallet()
	if len(wallet.Accounts) != 0 {
		t.Error("初始化钱包失败")
	}
	addr := wallet.AddAccount()
	if addr == "" {
		t.Error("添加账户失败")
	}
	addrs := wallet.GetAllAddr()
	if len(addrs) != 1 {
		t.Error("获取账户失败")
	}
}
