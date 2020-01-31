# dist-daemon-cosmos
Distributed Daemonize service on Cosmos SDK

``` 

go run cmd/main.go init bcc1 --chain-id bccchain --home ./chainroot1

go run cmd/main.go init bcc2 --chain-id bccchain --home ./chainroot2

go run cmd/main.go init bcc2 --chain-id bccchain --home ./chainroot3

```
``` 
go run cmd/main.go keys add operator --home ./chainroot1 
go run cmd/main.go add-genesis-account $(go run ./cmd/main.go keys show operator -a --home ./chainroot1) 1000nametoken,100000000stake  --home ./chainroot1


go run cmd/main.go  gentx --name operator  --home ./chainroot1  --home-client ./chainroot1
go run cmd/main.go  collect-gentxs --home ./chainroot1

go run cmd/main.go  validate-genesis --home ./chainroot1

go run cmd/main.go  tendermint show-node-id --home ./chainroot1

```

``` 
#overwrite chainroot2/config/genesis.json with first node's genesis.json

set allow_duplicate_ip = false in config/config.toml in all chainroot directories

change ports in chainroot2/config/config.toml
change ports in chainroot3/config/config.toml



and add [p2p] : persistent_peers = "d6809d7052471c7fc6c43c95ba97db597826ed29@127.0.0.1:26656"

```

``` 
go run cmd/main.go  start --home ./chainroot1

go run cmd/main.go  start --home ./chainroot2

go run cmd/main.go  start --home ./chainroot3

```
