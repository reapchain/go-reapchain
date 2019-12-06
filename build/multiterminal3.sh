#!/bin/bash
shopt -s nocasematch
osascript -e `tell application "Terminal"`  -e `do script  "/Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/bin/geth --networkid 2017 --port 30301 --nodiscover --datadir ./node1/data --mine --minerthreads 1 --syncmode "full"  --rpc --rpcport 8541 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" --unlock 0 --password /Users/yongilchoi/go/src/github.com/ethereum/go-ethereum/build/passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug"`  -e `end tell`
shopt -u nocasematch