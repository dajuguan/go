# 接口自动继承
```go
type readerWithCache struct {
	Reader // safe for concurrent read

	// Previously resolved state entries.
	accounts    map[common.Address]*types.StateAccount
	accountLock sync.RWMutex

	// List of storage buckets, each of which is thread-safe.
	// This reader is typically used in scenarios requiring concurrent
	// access to storage. Using multiple buckets helps mitigate
	// the overhead caused by locking.
	storageBuckets [16]struct {
		lock     sync.RWMutex
		storages map[common.Address]map[common.Hash]common.Hash
	}
}
```

### geth中交易池相关逻辑
```
func (p *BlobPool) SubscribeTransactions

// eth/handler.go 在产生新交易之后，会通过 p2p 网络广播出去
func (h *handler) txBroadcastLoop() {
	defer h.wg.Done()
	for {
		select {
		case event := <-h.txsCh: // 这里监听新交易信息
			h.BroadcastTransactions(event.Txs)
		case <-h.txsSub.Err():
			return
		}
	}
```