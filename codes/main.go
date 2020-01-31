package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/dajuguan/go/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("getBalance -address ADDRESS -获取ADDRESS的余额")
	fmt.Println(" createblockchain -address ADDRESS -创建挖创世区块账户地址")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT -FROM发送给TO一定AMOUNT数量的")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printBlocks() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Previos Hash %x\n", block.PrevHash)
		fmt.Printf("Data:%s\n", block.Transactions)
		fmt.Printf("Hash is %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}
func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished")
}
func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	utxos := chain.FindUTXO(address)
	for _, out := range utxos {
		balance += out.Value
	}
	fmt.Printf("Balance of address %s is: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Printf("[%s] Transfer %d To [%s] Successs", from, amount, to)
}
func (cli *CommandLine) Run() {
	cli.validateArgs()

	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createblockchain := flag.NewFlagSet("create", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "余额地址")
	createBlockchainAddr := createblockchain.String("address", "", "创建创世区块地址")
	sendFrom := sendCmd.String("from", "", "发送发地址")
	sendTo := sendCmd.String("to", "", "接收方地址")
	sendAmount := sendCmd.Int("amount", 0, "金额")
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
	case "create":
		err := createblockchain.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}
	if createblockchain.Parsed() {
		if *createBlockchainAddr == "" {
			createblockchain.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddr)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printBlocks()
	}

}

func main() {
	defer os.Exit(0)
	cli := CommandLine{}
	cli.Run()
}
