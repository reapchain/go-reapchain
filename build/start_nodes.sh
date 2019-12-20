./bin/geth --networkid 2017 --port 30501  --datadir ./node1/data --mine --minerthreads 1 \
--syncmode "full"  --rpc --rpcport 8541 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug


./bin/geth --networkid 2017 --port 30502  --datadir ./node2/data --mine --minerthreads 1 \
--syncmode "full"  --rpc --rpcport 8542 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug

./bin/geth --networkid 2017 --port 30503  --datadir ./node3/data --mine --minerthreads 1 \
--syncmode "full"  --rpc --rpcport 8543 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug

./bin/geth --networkid 2017 --port 30504  --datadir ./node4/data --mine --minerthreads 1 \
--syncmode "full"  --rpc --rpcport 8544 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.2 --ipcdisable --debug