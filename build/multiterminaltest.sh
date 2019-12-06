#!/bin/bash

echo -e "Open 30301 terminal\n\n"
open -na Terminal.app /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build -e "`/Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/bin/geth --networkid 2017 --port 30301 --nodiscover --datadir ./node1/data --mine --minerthreads 1 \
--syncmode "full"  --rpc --rpcport 8541 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" \
--unlock 0 --password /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug`"

echo -e "Open 30302 terminal\n\n"
open -na Terminal.app /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build -e "`/Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/bin/geth --networkid 2017 --port 30302 --nodiscover --datadir ./node2/data --mine --minerthreads 1 \
--syncmode "full"  --rpc --rpcport 8542 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" \
--unlock 0 --password /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug`"
echo -e "Open 30303 terminal\n\n"
open -na Terminal.app /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build -e "`/Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/bin/geth --networkid 2017 --port 30303 --nodiscover --datadir ./node3/data --mine --minerthreads 1 \
--syncmode "full"  --rpc --rpcport 8543 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" \
--unlock 0 --password /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug`"
echo -e "Open 30304 terminal\n\n"
open -na Terminal.app /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build  "`/Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/bin/geth --networkid 2017 --port 30304 --nodiscover --datadir ./node4/data --mine --minerthreads 1 \
--syncmode "full"  --rpc --rpcport 8544 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" \
--unlock 0 --password /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug`"