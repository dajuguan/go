package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"

	"github.com/dajuguan/go/blockchain"
	"github.com/vrecan/death"
)

const (
	protocol      = "tcp"
	version       = 1
	commandLength = 12
)

var (
	nodeAddress     string //对应节点端口
	minerAddress    string //对应矿工节点端口
	KnownNodes      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	memoryPool      = make(map[string]blockchain.Transaction)
)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block    []byte
}

type GetBlocks struct {
	AddrFrom string
}

type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

//多个区块
type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type Tx struct {
	AddrFrom    string
	Transaction []byte
}

//对于区块同步非常重要
type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func CmdToBytes(cmd string) []byte {
	var bytes [commandLength]byte
	for i, c := range cmd {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func BytesToCmd(bytes []byte) string {
	var cmd []byte
	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}
	return fmt.Sprintf("%s", cmd)
}

func CloseDB(chain *blockchain.BlockChain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}

func GobEncode(data interface{}) []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func ExtractCmd(req []byte) []byte {
	return req[:commandLength]
}

func SendAddr(address string) {
	nodes := Addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := GobEncode(nodes)
	req := append(CmdToBytes("addr"), payload...)
	SendData(address, req)
}

func SendInv(address, kind string, items [][]byte) {
	inventory := Inv{address, kind, items}
	payload := GobEncode(inventory)
	req := append(CmdToBytes("inv"), payload...)
	SendData(address, req)
}

func SendTx(addr string, tnx *blockchain.Transaction) {
	data := Tx{nodeAddress, tnx.Serialize()}
	payload := GobEncode(data)
	req := append(CmdToBytes("tx"), payload...)
	SendData(addr, req)
}

func SendVersion(addr string, chain *blockchain.BlockChain) {
	fmt.Println()
	bestHeight := chain.GetBestHeight()
	payload := GobEncode(Version{version, bestHeight, nodeAddress})
	req := append(CmdToBytes("version"), payload...)
	SendData(addr, req)
}

func SendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		updatedNodes := KnownNodes[:]
		for _, node := range KnownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}
		KnownNodes = updatedNodes
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func SendGetBlocks(address string) {
	payload := GobEncode(GetBlocks{nodeAddress})
	req := append(CmdToBytes("getblocks"), payload...)
	SendData(address, req)
}

func SendGetData(address, kind string, id []byte) {
	payload := GobEncode(GetData{nodeAddress, kind, id})
	req := append(CmdToBytes("getdata"), payload...)
	SendData(address, req)
}
func SendBlock(addr string, b *blockchain.Block) {
	data := Block{nodeAddress, b.Serialize()}
	payload := GobEncode(data)
	req := append(CmdToBytes("block"), payload...)
	SendData(addr, req)
}

//添加节点并同步区块
func HandleAddr(req []byte) {
	var buffer bytes.Buffer
	var payload Addr

	buffer.Write(req[commandLength:])
	dec := gob.NewDecoder(&buffer)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("共有 %d个节点", len(KnownNodes))
	RequestBlocks()
}

func RequestBlocks() {
	for _, node := range KnownNodes {
		SendGetBlocks(node)
	}
}

func HandleBlock(req []byte, chain *blockchain.BlockChain) {
	var buffer bytes.Buffer
	var payload Block

	buffer.Write(req[commandLength:])
	dec := gob.NewDecoder(&buffer)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blockData := payload.Block
	block := blockchain.Deserialize(blockData)

	fmt.Printf("Added Block %x\n", block.Hash)
	chain.AddBlock(block)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		SendGetData(payload.AddrFrom, "block", blockHash)
		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{chain}
		UTXOSet.Reindex()
	}
}

//发送所有区块的Hash
func HandleGetBlocks(req []byte, chain *blockchain.BlockChain) {
	var buffer bytes.Buffer
	var payload GetBlocks

	buffer.Write(req[commandLength:])
	dec := gob.NewDecoder(&buffer)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := chain.GetBlockHashes()
	SendInv(payload.AddrFrom, "block", blocks)
}

func HandleGetData(req []byte, chain *blockchain.BlockChain) {
	var buffer bytes.Buffer
	var payload GetData

	buffer.Write(req[commandLength:])
	dec := gob.NewDecoder(&buffer)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.Type == "block" {
		block, err := chain.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}
		SendBlock(payload.AddrFrom, &block)
	}
	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := memoryPool[txID]

		SendTx(payload.AddrFrom, &tx)
	}

}

func HandleVersion(req []byte, chain *blockchain.BlockChain) {
	var buffer bytes.Buffer
	var payload Version

	buffer.Write(req[commandLength:])
	dec := gob.NewDecoder(&buffer)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	bestHeight := chain.GetBestHeight()
	otherHeight := payload.BestHeight

	if bestHeight < otherHeight {
		SendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		SendVersion(payload.AddrFrom, chain)
	}
	if !NodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}

}

func HandleTx(req []byte, chain *blockchain.BlockChain) {
	var buffer bytes.Buffer
	var payload Tx

	buffer.Write(req[commandLength:])
	dec := gob.NewDecoder(&buffer)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)
	memoryPool[hex.EncodeToString(tx.ID)] = tx
	fmt.Printf("%s, %d", nodeAddress, len(memoryPool))

	//中心节点
	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				SendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(memoryPool) >= 2 && len(minerAddress) > 0 {
			MineTx(chain)
		}
	}

}

func MineTx(chain *blockchain.BlockChain) {
	var txs []*blockchain.Transaction
	for id := range memoryPool {
		fmt.Println("tx : %s \n", memoryPool[id].ID)
		tx := memoryPool[id]
		if chain.VerifyTransaction(&tx) {
			txs = append(txs, &tx)
		}
	}
	if len(txs) == 0 {
		fmt.Println("所有的交易无效")
		return
	}
	cbtx := blockchain.CoinbaseTx(minerAddress, "")
	txs = append(txs, cbtx)

	newBlock := chain.MineBlock(txs)
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()
	fmt.Println("New block mined")

	for _, tx := range txs {
		txID := hex.EncodeToString(tx.ID)
		delete(memoryPool, txID)
	}
	for _, node := range KnownNodes {
		if node != nodeAddress {
			SendInv(node, "block", [][]byte{newBlock.Hash})
		}
	}
	if len(memoryPool) > 0 {
		MineTx(chain)
	}
}

func HandleInv(req []byte, chain *blockchain.BlockChain) {
	var buffer bytes.Buffer
	var payload Inv

	buffer.Write(req[commandLength:])
	dec := gob.NewDecoder(&buffer)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.Type == "block" {
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]
		SendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]
		if memoryPool[hex.EncodeToString(txID)].ID == nil {
			SendGetData(payload.AddrFrom, "tx", txID)
		}
	}

}

func NodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}
	return false
}

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	req, err := ioutil.ReadAll(conn)
	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}
	cmd := BytesToCmd(req[:commandLength])
	fmt.Println("get command:", cmd)
	switch cmd {
	case "addr":
		HandleAddr(req)
	case "block":
		HandleBlock(req, chain)
	case "inv":
		HandleInv(req, chain)
	case "getblocks":
		HandleGetBlocks(req, chain)
	case "getdata":
		HandleGetData(req, chain)
	case "tx":
		HandleTx(req, chain)
	case "version":
		HandleVersion(req, chain)
	default:
		fmt.Printf("Unknown Command:%s", cmd)
	}
}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	minerAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Database.Close()
	//冗余操作，确保正确
	go CloseDB(chain)

	if nodeAddress != KnownNodes[0] {
		SendVersion(KnownNodes[0], chain)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go HandleConnection(conn, chain)
	}
}
