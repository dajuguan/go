package blockchain

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/dgraph-io/badger"
)

type UTXOSet struct {
	BC *BlockChain
}

var (
	utxoPrefix   = []byte("utxo-")
	prefixLength = len(utxoPrefix)
)

//查找账户的UTXO
func (u UTXOSet) FindUnspentTransactions(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	db := u.BC.Database
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item()
			var v []byte
			if err := item.Value(func(val []byte) error {
				v = val
				return nil
			}); err != nil {
				return err
			}
			outs := DeserializeOutputs(v)
			for _, out := range outs.Outputs {
				if out.isLockedWithPubKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}

			}
		}
		return nil
	})
	Handle(err)
	return UTXOs
}

func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	utxos := make(map[string][]int)
	total := 0
	db := u.BC.Database

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item()
			key := item.Key()
			k := bytes.TrimPrefix(key, utxoPrefix)
			txID := hex.EncodeToString(k)
			var v []byte
			if err := item.Value(func(val []byte) error {
				v = val
				return nil
			}); err != nil {
				return err
			}
			outs := DeserializeOutputs(v)
			for outIdx, out := range outs.Outputs {
				if out.isLockedWithPubKey(pubKeyHash) && total < amount {
					total += out.Value
					utxos[txID] = append(utxos[txID], outIdx)
					// if total >= amount {
					// 	break Work
					// }
				}
			}
		}
		return nil
	})
	Handle(err)

	return total, utxos
}

// func (chain *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
// 	var unspentTxs []Transaction
// 	spentTXOs := make(map[string][]int)
// 	iter := chain.Iterator()

// 	for {
// 		block := iter.Next()
// 		for _, tx := range block.Transactions {
// 			txID := hex.EncodeToString(tx.ID)
// 			//查询交易单的所有输出，记录未花费的交易输出
// 		Outputs:
// 			for outIdx, out := range tx.Outputs {
// 				if spentTXOs[txID] != nil {
// 					for _, spentOut := range spentTXOs[txID] {
// 						//交易单里面有一个输出的ID对应上，那么就不用查询了
// 						if spentOut == outIdx {
// 							continue Outputs
// 						}
// 					}
// 				}
// 				if out.isLockedWithPubKey(pubKeyHash) {
// 					unspentTxs = append(unspentTxs, *tx)
// 				}
// 			}

// 			//记录花费过的输入
// 			if tx.IsCoinbase() == false {
// 				for _, in := range tx.Inputs {
// 					if in.UserKey(pubKeyHash) {
// 						inTxID := hex.EncodeToString(in.ID)
// 						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
// 					}
// 				}
// 			}
// 		}
// 		if len(block.PrevHash) == 0 {
// 			break
// 		}
// 	}

// 	return unspentTxs
// }

func (u *UTXOSet) Reindex() {
	db := u.BC.Database
	u.DeleteByPrefix(utxoPrefix)

	UTXO := u.BC.FindUTXO()
	err := db.Update(func(txn *badger.Txn) error {
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			key = append(utxoPrefix, key...)
			err = txn.Set(key, outs.Serialize())
			Handle(err)
		}
		return nil
	})
	Handle(err)
}

//加入最新的区块，更新UTXO
//1.更新过去的UTXO
//2.加入新的交易输出：包含coinbase和其他的输出
func (u *UTXOSet) Update(block *Block) {
	db := u.BC.Database
	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range block.Transactions {
			for _, in := range tx.Inputs {
				inID := append(utxoPrefix, in.ID...)
				item, err := txn.Get(inID)
				if err != nil {
					return err
				}
				var v []byte
				err = item.Value(func(val []byte) error {
					v = val
					return nil
				})
				outs := DeserializeOutputs(v)
				//更新过去inID相关的输出outs，只要out的ID与当前交易的ID不同就是未引用
				outputs := TxOutputs{}
				for outID, out := range outs.Outputs {
					if outID != in.Out {
						outputs.Outputs = append(outputs.Outputs, out)
					}
				}
				//如果没有out，删除;否则更新
				if len(outputs.Outputs) == 0 {
					if err := txn.Delete(inID); err != nil {
						return err
					}
				} else {
					if err := txn.Set(inID, outputs.Serialize()); err != nil {
						return err
					}
				}

			}
			//插入新的UTXO：coinbase交易输出，普通交易输出
			newOutputs := TxOutputs{}
			for _, out := range tx.Outputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}
			txID := append(utxoPrefix, tx.ID...)
			txn.Set(txID, newOutputs.Serialize())
			return nil

		}
		return nil
	})
	Handle(err)
}

//计算所有UTXO中的交易数目
func (u UTXOSet) CountTransactions() int {
	db := u.BC.Database
	counter := 0
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			counter++
		}
		return nil
	})
	Handle(err)
	return counter
}

func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
	deleteKeys := func(keysForDelete [][]byte) error {
		err := u.BC.Database.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		})
		return err
	}

	collectSize := 100000

	u.BC.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					log.Panic(err)
				}
				keysForDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				log.Panic(err)
			}
		}
		return nil

	})

}
