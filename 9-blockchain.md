# 开发区块链

```
go run main.go create -address "alice"
go run main.go getbalance -address alice
go run main.go send -from "alice" -to "bob" -amount 50
go run main.go getbalance -address bob
```

钱包
```
go run main.go createwallet
go run main.go listaddrs
```

```
go run main.go createwallet
go run main.go createwallet
go run main.go createbc -address 1Jm6E8spEtZrGc8KmNbstmd7LviMjuQHaA
go run main.go listaddr
go run main.go getbalance -address 17cteVpL7ZqhXwrr8ZGRVoN67HHHtxFBMU
go run main.go send -to 1FH4gDNzKYnkKis5v4kKijc2arWftTmL8 -from 17cteVpL7ZqhXwrr8ZGRVoN67HHHtxFBMU -amount 50
go run main.go getbalance -address 1FH4gDNzKYnkKis5v4kKijc2arWftTmL8
```

UTXO重建
```
go run main.go create -address 
go run main.go createbc -address 17cteVpL7ZqhXwrr8ZGRVoN67HHHtxFBMU
go run main.go getbalance -addr 17cteVpL7ZqhXwrr8ZGRVoN67HHHtxFBMU
# 必须先reindex才可以
go run main.go reindexutxo
go run main.go getbalance -address 17cteVpL7ZqhXwrr8ZGRVoN67HHHtxFBMU
```