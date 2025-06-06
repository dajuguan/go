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