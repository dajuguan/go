package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/dajuguan/go/codes/blockchain"
	"github.com/dajuguan/go/codes/network"
	"github.com/dajuguan/go/codes/wallet"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("getBalance -address ADDRESS -获取ADDRESS的余额")
	fmt.Println("createbc -address ADDRESS -创建挖创世区块账户地址")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT  -mine -FROM发送给TO一定AMOUNT数量的")
	fmt.Println("createwallet -创建新的钱包")
	fmt.Println("listaddrs -列出钱包所有地址")
	fmt.Println("reindexutxo -重建UTXO集")
	fmt.Println("startnode -miner ADDRESS -启动ID为ADDRESS的节点")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) StartNode(nodeID, minerAddress string) {
	fmt.Printf("启动节点%s\n", nodeID)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("开始挖矿，接收奖励的地址为:", minerAddress)
		} else {
			log.Panic("错误的挖矿地址")
		}
	}
	network.StartServer(nodeID, minerAddress)
}

func (cli *CommandLine) reindexUTXO(nodeID string) {
	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()
	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set", count)
}

func (cli *CommandLine) printBlocks(nodeID string) {
	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		// fmt.Printf("Previos Hash %x\n", block.PrevHash)
		// fmt.Printf("Data:%s\n", block.Transactions)
		// fmt.Printf("Hash is %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}
func (cli *CommandLine) createBlockChain(address, nodeID string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}
	print("building...")
	chain := blockchain.InitBlockChain(address, nodeID)
	chain.Database.Close()

	//初始化UTXO数据集
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()
	fmt.Println("Finished")
}
func (cli *CommandLine) getBalance(address string, nodeID string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}
	chain := blockchain.ContinueBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	balance := 0

	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	utxos := UTXOSet.FindUnspentTransactions(pubKeyHash)
	for _, out := range utxos {
		balance += out.Value
	}
	fmt.Printf("Balance of address %s is: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int, nodeID string, mineNow bool) {
	if !wallet.ValidateAddress(from) {
		log.Panic("Address is not Valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("Address is not Valid")
	}
	chain := blockchain.ContinueBlockChain(nodeID)
	utxoset := &blockchain.UTXOSet{chain}

	defer chain.Database.Close()

	//创建钱包
	wallets, err := wallet.CreateWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := blockchain.NewTransaction(&wallet, to, amount, utxoset)
	fmt.Println("===============")

	//本机挖矿，还是远程挖矿
	if mineNow {
		cbtx := blockchain.CoinbaseTx(from, "")
		txs := []*blockchain.Transaction{cbtx, tx}
		block := chain.MineBlock(txs)
		utxoset.Update(block)
	} else {
		network.SendTx(network.KnownNodes[0], tx)
		fmt.Println("发送交易")
	}
	utxoset.Reindex()
	fmt.Printf("[%s] Transfer %d To [%s] Successs", from, amount, to)
}

func (cli *CommandLine) listWalletsAddr(nodeID string) {
	wallets, err := wallet.CreateWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAllAddresses()
	for _, addr := range addresses {
		fmt.Println(addr)
	}
}

func (cli *CommandLine) createWallet(nodeID string) {
	wallets, err := wallet.CreateWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	address := wallets.AddWallet()
	wallets.SaveFile(nodeID)
	fmt.Println("新创建的钱包地址为:", address)

}

func (cli *CommandLine) Run() {
	cli.validateArgs()
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env is not set!")
		runtime.Goexit()
	}

	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createblockchain := flag.NewFlagSet("createbc", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressedCmd := flag.NewFlagSet("listaddrs", flag.ExitOnError)
	reindexutxoCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "余额地址")
	createBlockchainAddr := createblockchain.String("address", "", "创建创世区块地址")
	sendFrom := sendCmd.String("from", "", "发送发地址")
	sendTo := sendCmd.String("to", "", "接收方地址")
	sendAmount := sendCmd.Int("amount", 0, "金额")
	sendMine := sendCmd.Bool("mine", false, "在相同的节点上同时挖矿")
	startNodeMiner := startNodeCmd.String("miner", "", "开启挖矿")

	switch os.Args[1] {
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createbc":
		err := createblockchain.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createwallet":
		createWalletCmd.Parse(os.Args[2:])
	case "listaddrs":
		listAddressedCmd.Parse(os.Args[2:])
	case "reindexutxo":
		reindexutxoCmd.Parse(os.Args[2:])
	case "startnode":
		startNodeCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		runtime.Goexit()
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress, nodeID)
	}
	if createblockchain.Parsed() {
		if *createBlockchainAddr == "" {
			createblockchain.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddr, nodeID)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)
	}

	if printChainCmd.Parsed() {
		cli.printBlocks(nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}
	if listAddressedCmd.Parsed() {
		cli.listWalletsAddr(nodeID)
	}

	if reindexutxoCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}
	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			runtime.Goexit()
		}
		cli.StartNode(nodeID, *startNodeMiner)
	}

}
