#!/bin/bash
## made by yichoi 2011-11-22 for automatic booting for the geth of istanbul base source


# initialize geth 4 nodes

# pwd = $PROJECTDIR/build/

# 아래 명령은 한번만 수행하고, 그 다음부터는 # mark 해서 실행 안되게 할것.

# 4개는 3+1 = validator, 1개는 go land에서 사용하는 디버깅용.
#../istanbul-tools/build/bin/istanbul setup --num 5 --nodes  --save --verbose > out2.txt

lsof -i:30501-30505

rm -rf node1 node2 node3 node4 node5

mkdir node1 node2 node3 node4 node5

ls -l node*





ls -l

cat ~/go/src/github.com/ethereum/go-ethereum/build/static-nodes.json

echo -e "update the IP and port : same as below

[
    "enode://dd333ec28f0a8910c92eb4d336461eea1c20803eed9cf2c056557f986e720f8e693605bba2f4e8f289b1162e5ac7c80c914c7178130711e393ca76abc1d92f57@192.168.0.2:30303?discport=30301", \n
    "enode://1bb6be462f27e56f901c3fcb2d53a9273565f48e5d354c08f0c044405b29291b405b9f5aa027f3a75f9b058cb43e2f54719f15316979a0e5a2b760fff4631998@192.168.0.2:30304?discport=30301", \n
    "enode://0df02e94a3befc0683780d898119d3b675e5942c1a2f9ad47d35b4e6ccaf395cd71ec089fcf1d616748bf9871f91e5e3d29c1cf6f8f81de1b279082a104f619d@192.168.0.2:30305?discport=30301", \n
    "enode://3fe0ff0dd2730eaac7b6b379bdb51215b5831f4f48fa54a24a0298ad5ba8c2a332442948d53f4cd4fd28f373089a35e806ef722eb045659910f96a1278120516@192.168.0.2:30306?discport=30301"  \n

]
"



cd ~/go/src/github.com/ethereum/go-ethereum/build

mkdir -p node1/data/geth
mkdir -p node2/data/geth
mkdir -p node3/data/geth
mkdir -p node4/data/geth
mkdir -p node5/data/geth

ifconfig en0



# 새 계정 4개를 각 노드마다 새로 생성한다. - passwd.txt의 암호 파일에서 읽어와서..
#./bin/geth --datadir node1/data account new --password passwd.txt
#./bin/geth --datadir node1/data account import --password passwd.txt ./node1/data/geth/nodekey
# Address: {9480652c2be98237715e38d0e72367088f84233e}
#./bin/geth --datadir node2/data account new --password passwd.txt
#./bin/geth --datadir node1/data account import --password passwd.txt ./node2/data/geth/nodekey
# Address: {1d423d94bdea465eb5260b2a61f2c02ab66bd89e}
#./bin/geth --datadir node3/data account new --password passwd.txt
#Address: {9cca27c05d5ecc2bf4ba44cdb336ffbf4fd3d0e3}
#./bin/geth --datadir node4/data account new --password passwd.txt
#Address: {21d3c0a6c14f934ca420976f1ed2d2e2478b83c6}
#./bin/geth --datadir node5/data account new --password passwd.txt
#Address: {520d72fcca400284690adff119a93dbaf3f24133}


## 총 5개 어카운트를 생성한다.

echo -e "To add accounts to the initial block, edit the genesis.json file in the lead node’s working directory and
update the alloc field with the account(s) that were generated at previous step"

cd ~/go/src/github.com/ethereum/go-ethereum/build

cp genesis.json node1
cp genesis.json node2
cp genesis.json node3
cp genesis.json node4
cp genesis.json node5

cp static-nodes.json node1/data/
cp static-nodes.json node2/data/
cp static-nodes.json node3/data/
cp static-nodes.json node4/data/
cp static-nodes.json node5/data/

cd ~/go/src/github.com/ethereum/go-ethereum/build
cp ./0/nodekey node1/data/geth
cp ./1/nodekey node2/data/geth
cp ./2/nodekey node3/data/geth
cp ./3/nodekey node4/data/geth
cp ./4/nodekey node5/data/geth

./bin/geth --datadir node1/data account import --password passwd.txt ./node1/data/geth/nodekey
./bin/geth --datadir node2/data account import --password passwd.txt ./node2/data/geth/nodekey
./bin/geth --datadir node3/data account import --password passwd.txt ./node3/data/geth/nodekey
./bin/geth --datadir node4/data account import --password passwd.txt ./node4/data/geth/nodekey
./bin/geth --datadir node5/data account import --password passwd.txt ./node5/data/geth/nodekey

pwd

cd ~/go/src/github.com/ethereum/go-ethereum/build

# 위에서 계정을 만든 것으로 각 노드의 genesis.json의 alloc에 이더를 할당해주어야함.
#"alloc": {
#        "8aaa61e122c0646cea852bcc2b92b9af39111498": {
#            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
#        },
#        "985f84e7cd9c1d980b512e771b86234663a561ff": {
#            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
#        },
#        "a57c86048931e6d2501ac31a19de3bd41fba5467": {
#            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
#        },
#        "bd72ddeecafd9afb45173b69b8b3645c6dd83852": {
#            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
#        }

./bin/geth --datadir ./node1/data init ./node1/genesis.json
./bin/geth --datadir ./node2/data init ./node2/genesis.json
./bin/geth --datadir ./node3/data init ./node3/genesis.json
./bin/geth --datadir ./node4/data init ./node4/genesis.json
./bin/geth --datadir ./node5/data init ./node5/genesis.json
# ./bin/geth --datadir=./node1 --port 5000 --rpcport 6000 –-ipcdisable --mine --minerthreads 1 --syncmode "full"

ifconfig | grep netmask|awk '{print $2}'


#합의엔진 테스트환경 설정

#--verbosity 옵션은 로그 출력 수준을 지정할 수 있습니다.

# (0=silent, 1=error, 2=warn, 3=info, 4=core, 5=debug, 6=detail)

# RUN 3 개,, 디버그 1개

#for validators :

# Validator 1

# ./bin/geth --datadir "./node1/data" --mine --minerthreads 1 --syncmode "full"

# forground service 를 띄우면서 console 접속함

cd ~/go/src/github.com/ethereum/go-ethereum/build

./bin/geth --networkid 2017 --port 30501  --datadir ./node1/data --mine --minerthreads 1 --syncmode "full" \
--rpc --rpcport 8545 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" --ipcdisable --debug \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.100 console 2>> ./node1/geth.log

./bin/geth --networkid 2017 --port 30502  --datadir ./node2/data --mine --minerthreads 1 --syncmode "full" \
--rpc --rpcport 8546 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" --ipcdisable --debug \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.100 console 2>> ./node2/geth.log

./bin/geth --networkid 2017 --port 30503  --datadir ./node3/data --mine --minerthreads 1 --syncmode "full" \
--rpc --rpcport 8547 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" --ipcdisable --debug \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.100 console 2>> ./node3/geth.log

./bin/geth --networkid 2017 --port 30504  --datadir ./node4/data --mine --minerthreads 1 --syncmode "full" \
--rpc --rpcport 8548 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" --ipcdisable --debug \
--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.100 console 2>> ./node4/geth.log







#./bin/geth --networkid 2017 --port 30303 --nodiscover --datadir ./node4/data --mine --minerthreads 1 --syncmode "full" \
#--rpc --rpcport 8548 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug"  --ipcdisable\
#--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.2

# reqular nodes & Debug node in Go land
# Go land의 Run-> Edit configuraiton -> Parameter argument에 아래 넣어서,
# Goland의 디버그를 샐행새켜서 볼것,
#./bin/geth --networkid 2017 --datadir ./node5/data --port 30304 --rpcport 8549 --rpccorsdomain "*" \
#--rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug"  \
#--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.2 --dev console






#See if the any geth nodes are running.
ps | grep geth
##Kill geth processes
##killall -INT geth

# $ chmod +x startistanbul.sh
# $ ./startistanbul.sh
# $ ps


lsof -i:30301-30305

