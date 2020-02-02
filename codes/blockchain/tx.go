package blockchain

type TxOutput struct {
	Value  int
	Pubkey string
}

type TxInput struct {
	ID  []byte //关键的交易
	Out int    //关联输出的index
	Sig string
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.Pubkey == data
}
