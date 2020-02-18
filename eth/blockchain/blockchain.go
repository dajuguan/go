package blockchain

type State struct {
}

type BlockChain struct {
	Blocks []*Block
	State  State
}

const BENEFICIARY = 100

//没有父Hash的Block
func (chain *BlockChain) AddBlock(b *Block) {
	lastBlock := chain.Blocks[len(chain.Blocks)-1]
	state := chain.State
	if Validate(lastBlock, b, state) {
		//增加Coinbase交易
		cb := CoinbaseTx([]byte("1111"), BENEFICIARY)
		b.Transactions = append(b.Transactions, cb)
		chain.Blocks = append(chain.Blocks, b)
	}
}

func InitBlockChain() *BlockChain {
	var bc BlockChain
	genesis := Genesis()
	bc = BlockChain{[]*Block{genesis}, State{}}
	return &bc
}
