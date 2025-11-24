# test golang's scheduler issue
```
cd hash_bench
# watch cpu usage with htop
# 1 thread: ~ 150k
taskset -c 0 go run main.go -n 100000000 -t 1 - r 1000
# 8 thread: ~80k, not linear with num of threads
taskset -c 0-7 go run main.go -n 100000000 -t 8 - r 1000

```

## Ref
https://github.com/ethstorage/zk-decoder/blob/main/golang/cmd/hash_bench/main.go