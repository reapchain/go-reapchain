## made by yichoi 2011-11-22 for automatic booting for the geth of istanbul base source


# initialize geth 4 nodes

# pwd = $PROJECTDIR/build/

rm -rf node1 node2 node3 node4

mkdir node1 node2 node3 node4

ls -l node*

cd node1

# ../../istanbul-tools/build/bin/istanbul setup --num 4 --nodes  --save --verbose

ls -l

cat static-nodes.json

echo -e "update the IP and port : same as below

[
    "enode://dd333ec28f0a8910c92eb4d336461eea1c20803eed9cf2c056557f986e720f8e693605bba2f4e8f289b1162e5ac7c80c914c7178130711e393ca76abc1d92f57@127.0.0.1:30301?discport=0", \n
    "enode://1bb6be462f27e56f901c3fcb2d53a9273565f48e5d354c08f0c044405b29291b405b9f5aa027f3a75f9b058cb43e2f54719f15316979a0e5a2b760fff4631998@127.0.0.1:30302?discport=0", \n
    "enode://0df02e94a3befc0683780d898119d3b675e5942c1a2f9ad47d35b4e6ccaf395cd71ec089fcf1d616748bf9871f91e5e3d29c1cf6f8f81de1b279082a104f619d@127.0.0.1:30303?discport=0", \n
    "enode://3fe0ff0dd2730eaac7b6b379bdb51215b5831f4f48fa54a24a0298ad5ba8c2a332442948d53f4cd4fd28f373089a35e806ef722eb045659910f96a1278120516@127.0.0.1:30304?discport=0"  \n

]
"
cd ..

mkdir -p node1/data/geth
mkdir -p node2/data/geth
mkdir -p node3/data/geth
mkdir -p node4/data/geth

./bin/geth --datadir node1/data account new
##
#Your new account is locked with a password. Please give a password. Do not forget this password.
#Passphrase:
#Repeat passphrase:
#Address: {48417549913e78f04a18376ab51325d19d9c3739}
##

./bin/geth --datadir node2/data account new
./bin/geth --datadir node3/data account new
./bin/geth --datadir node4/data account new

echo -e "To add accounts to the initial block, edit the genesis.json file in the lead node’s working directory and
update the alloc field with the account(s) that were generated at previous step"

cp node1/genesis.json node2
cp node1/genesis.json node3
cp node1/genesis.json node4

cp node1/static-nodes.json node1/data/
cp node1/static-nodes.json node2/data/
cp node1/static-nodes.json node3/data/
cp node1/static-nodes.json node4/data/

cp node1/0/nodekey node1/data/geth
cp node1/1/nodekey node2/data/geth
cp node1/2/nodekey node3/data/geth
cp node1/3/nodekey node4/data/geth

pwd

cd node1
../bin/geth --datadir data init genesis.json

cd ..
cd node2
../bin/geth --datadir data init genesis.json

cd ..
cd node3
../bin/geth --datadir data init genesis.json

cd ..
cd node4
../bin/geth --datadir data init genesis.json

## 11..
#!/bin/bash

pwd

cd node1
nohup ../bin/geth --datadir data --nodiscover --istanbul.blockperiod 5 --syncmode full --mine --minerthreads 1 --verbosity 5 --networkid 2019 --rpc --rpcaddr 0.0.0.0 --rpcport 22000 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,istanbul --emitcheckpoints --port 30300 2>>node.log &
../bin/geth --datadir data --nodiscover --istanbul.blockperiod 5 --syncmode full --mine --minerthreads 1 --verbosity 5 --networkid 2019 --rpc --rpcaddr 0.0.0.0 --rpcport 22000 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,istanbul --emitcheckpoints --port 30300

cd ../node1
nohup geth --datadir data --nodiscover --istanbul.blockperiod 5 --syncmode full --mine --minerthreads 1 --verbosity 5 --networkid 2019 --rpc --rpcaddr 0.0.0.0 --rpcport 22001 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,istanbul --emitcheckpoints --port 30301 2>>node.log &

cd ../node2
nohup geth --datadir data --nodiscover --istanbul.blockperiod 5 --syncmode full --mine --minerthreads 1 --verbosity 5 --networkid 2019 --rpc --rpcaddr 0.0.0.0 --rpcport 22002 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,istanbul --emitcheckpoints --port 30302 2>>node.log &

cd ../node3
nohup geth --datadir data --nodiscover --istanbul.blockperiod 5 --syncmode full --mine --minerthreads 1 --verbosity 5 --networkid 2019 --rpc --rpcaddr 0.0.0.0 --rpcport 22003 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,istanbul --emitcheckpoints --port 30303 2>>node.log &

cd ../node4
nohup geth --datadir data --nodiscover --istanbul.blockperiod 5 --syncmode full --mine --minerthreads 1 --verbosity 5 --networkid 2019 --rpc --rpcaddr 0.0.0.0 --rpcport 22004 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,istanbul --emitcheckpoints --port 30304 2>>node.log &

#See if the any geth nodes are running.
ps | grep geth
Kill geth processes
killall -INT geth

# $ chmod +x startistanbul.sh
# $ ./startistanbul.sh
# $ ps



