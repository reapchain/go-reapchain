#!/bin/bash
osascript -e 'tell app "Terminal"
do script "./bin/geth --networkid 2017 --port 30301 --nodiscover --datadir ./node2/data --mine --minerthreads 1 --syncmode "full" \
 --rpc --rpcport 8546 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug"  \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug"
end tell'