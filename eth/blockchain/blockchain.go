package blockchain

type State struct {
}

type BlockChain struct {
	Blocks []*Block
	State  State
}

//没有父Hash的Block
func (chain *BlockChain) AddBlock(b *Block) {
	lastBlock := chain.Blocks[len(chain.Blocks)-1]
	state := chain.State
	if Validate(lastBlock, b, state) {
		chain.Blocks = append(chain.Blocks, b)
	}
}

func InitBlockChain() *BlockChain {
	var bc BlockChain
	genesis := Genesis()
	bc = BlockChain{[]*Block{genesis}, State{}}
	return &bc
}
