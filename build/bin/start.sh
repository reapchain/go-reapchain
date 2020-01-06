#!/bin/bash
## made by yichoi 2011-11-22 for automatic booting for the geth of istanbul base source


# initialize geth 4 nodes

# pwd = $PROJECTDIR/build/

# 아래 명령은 한번만 수행하고, 그 다음부터는 # mark 해서 실행 안되게 할것.

#./istanbul setup --num 7 --nodes  --save --verbose > recent.txt

# 중요 한번 istanbul-tools로 돌려서, 노드를 생성하면, 노드키가 생성됨으로,
# 그 다음 기동 시키려면, 반드시 Import .... nodekey로 해야, staitc_nodes.json, genesis.json을 재 편집할 필요가 없어진다.
# istanbul setup를 돌리면, genesis.json 파일을 생성해버려서, 여기서 Validator node 들의 20바이트 계정 정보가 alloc에 할당된다.
# 주의 사항: 모든 노드가 기동할때 똑같은 genesis.json과 networkid 를 일치시켜주어야 하고,
# geth port도 같아야 한다. 다만,,
# 각 노드들의 머신이 틀리때는 geth port, RPC port를 동일하게 주어야 통신됨.
# 한 머신 내에서 여러개의 geth node를 실행시 --port, --rpcport --nat pmp:127.0.0.1
# 여기서 --port, rpcport는 반드시 다르게 준다. 이 번호로 각 노드들을 구분하면, 이해하기 쉽다.
# node#1  ...#2... #6.... port 30501,, 30502..... 등.. 포트의 끝자리를 노드 폴더이름과 일치시켜줌,,


rm -rf geth*.log

rm -rf qman node1 node2 node3 node4 node5 node6



mkdir qman node1 node2 node3 node4 node5 node6

ls -l qman node*




mkdir -p qman/data/geth
mkdir -p node2/data/geth
mkdir -p node3/data/geth
mkdir -p node4/data/geth
mkdir -p node5/data/geth
mkdir -p node6/data/geth
mkdir -p node7/data/geth


# ifconfig  enp3s0
cp genesis.json qman
cp genesis.json node2
cp genesis.json node3
cp genesis.json node4
cp genesis.json node5
cp genesis.json node6
cp genesis.json node7


cp static-nodes.json qman/data/
cp static-nodes.json node2/data/
cp static-nodes.json node3/data/
cp static-nodes.json node4/data/
cp static-nodes.json node5/data/
cp static-nodes.json node6/data/
cp static-nodes.json node7/data/


cp qmanager-nodes.json qman/data/
cp qmanager-nodes.json node1/data/
cp qmanager-nodes.json node2/data/
cp qmanager-nodes.json node3/data/
cp qmanager-nodes.json node4/data/
cp qmanager-nodes.json node5/data/
cp qmanager-nodes.json node6/data/
cp qmanager-nodes.json node7/data/


cp ./0/nodekey qman/data/geth
cp ./1/nodekey node1/data/geth
cp ./2/nodekey node2/data/geth
cp ./3/nodekey node3/data/geth
cp ./4/nodekey node4/data/geth
cp ./5/nodekey node5/data/geth
cp ./6/nodekey node6/data/geth

# 새 계정 4개를 각 노드마다 새로 생성한다. - passwd.txt의 암호 파일에서 읽어와서..
./geth --datadir qman/data account new --password passwd.txt
./geth --datadir node1/data account new --password passwd.txt
./geth --datadir node2/data account new --password passwd.txt
./geth --datadir node3/data account new --password passwd.txt
./geth --datadir node4/data account new --password passwd.txt
./geth --datadir node5/data account new --password passwd.txt
./geth --datadir node6/data account new --password passwd.txt


# 한번 istanbul setup를 실행시키면, 매번 노드키가 바뀜. 따라서,
# 아래처럼,, 두번째 이 쉘을 실행시 반드시 account new가 아닌 account import ... ./.../../nodekey 형태로 써주어야함.


#./geth --datadir qman/data account import  --password passwd.txt ./qman/data/geth/nodekey
#./geth --datadir node2/data account import  --password passwd.txt ./node2/data/geth/nodekey
#./geth --datadir node3/data account import  --password passwd.txt ./node3/data/geth/nodekey
#./geth --datadir node4/data account import  --password passwd.txt ./node4/data/geth/nodekey
#./geth --datadir node5/data account import  --password passwd.txt ./node5/data/geth/nodekey
#./geth --datadir node6/data account import  --password passwd.txt ./node6/data/geth/nodekey
#./geth --datadir node7/data account import  --password passwd.txt ./node7/data/geth/nodekey


# syntax : geth account import --password <passwordfile> <keyfile>



## 총 5개 어카운트를 생성한다.

echo -e "To add accounts to the initial block, edit the genesis.json file in the lead node’s working directory and
update the alloc field with the account(s) that were generated at previous step"

pwd

./geth --datadir ./qman/data init ./qman/genesis.json
./geth --datadir ./node1/data init ./node1/genesis.json
./geth --datadir ./node2/data init ./node2/genesis.json
./geth --datadir ./node3/data init ./node3/genesis.json
./geth --datadir ./node4/data init ./node4/genesis.json
./geth --datadir ./node5/data init ./node5/genesis.json
./geth --datadir ./node6/data init ./node6/genesis.json


#합의엔진 테스트환경 설정

#--verbosity 옵션은 로그 출력 수준을 지정할 수 있습니다.

# (0=silent, 1=error, 2=warn, 3=info, 4=core, 5=debug, 6=detail)

# RUN 3 개,, 디버그 1개

#for validators :

# Validator 1

# ./bin/geth --datadir "./node1/data" --mine --minerthreads 1 --syncmode "full"

# forground service 를 띄우면서 console 접속함

#./bin/geth --networkid 2017 --port 30501 --nodiscover --datadir ./node1/data --mine --minerthreads 1 --syncmode "full" \
#--rpc --rpcport 8545 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug" --debug \
#--unlock 0 --password passwd.txt --verbosity 6 --nat pmp:127.0.0.1
#console 2>> ./node1/geth.log

# 콘솔 없이 실행하기
#osascript
#tell application "Terminal"
#do script

## ============================================================================================================
# 부트노드를 리눅스에 ( 192.168.0.100:30301) 먼저 띄우고.

# 일반노드를 구동
# nohup ./bin/geth  --networkid 2017 --port 30506  --datadir ./node1/data  \
# --ipcdisable --rpc --rpcport 8546 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,miner,admin,debug"   \
# --unlock 0 --password passwd.txt --verbosity 6 --nat pmp:127.0.0.1 console 2>> ./geth1.log &

# ======마이너 노드 구동=======
# Validator 1/5
# opened udp port check, if opened, succeeded! will displayed
nc -uz 192.168.0.2 30501-30510

echo "Qman debug node in Goland started"
./startQman.sh

echo "Node1 started"
./startnode1.sh

echo "Node3 started"
./startnode2.sh

echo "Node4 started"
./startnode3.sh

echo "Node4 started"
./startnode4.sh

echo "Node4 started"
./startnode5.sh

echo "Node4 started"
./startnode6.sh

# Validator 5/5 or Debug and Monitoring node
#nohup ./bin/geth  --networkid 2017 --port 30505  --datadir ./node5/data --mine --minerthreads 1 --syncmode "full" \
#--ipcdisable --rpc --rpcport 8545 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,miner,admin,debug"   \
# --unlock 0 --password passwd.txt --verbosity 9 --nat pmp:127.0.0.1 console 2>> ./geth5.log &


# mist node in dapp : mist dapp를 구동 시키기 위해서는 반드시 ipcenable 해줘야 해서, --ipcdisable 뺀다.

#./bin/geth2  --networkid 2017 --port 30510 --datadir ./node6/data --mine --minerthreads 1 --syncmode "full" \
# --rpc --rpcport 8560 --rpccorsdomain "*" --rpcapi="db,eth,net,web3,personal,miner,admin,debug"   \
#--unlock 0 --password passwd.txt --verbosity 6 --nat pmp:127.0.0.1

# reqular nodes & Debug node in Go land
# Go land의 Run-> Edit configuraiton -> Parameter argument에 아래 넣어서,
# Goland의 디버그를 샐행새켜서 볼것,
#./bin/geth --networkid 2017 --datadir ./node5/data --port 30307 --rpcport 8549 --rpccorsdomain "*" \
#--rpcapi="db,eth,net,web3,personal,web3,miner,admin,debug"  \
#--unlock 0 --password passwd.txt --verbosity 6 --nat extip:192.168.0.2 console 2>> ./node5/geth.log


#See if the any geth nodes are running.
ps -ef | grep geth >> background_processid.txt
##Kill geth processes
##killall -INT geth

# $ chmod +x startistanbul.sh
# $ ./startistanbul.sh
# $ ps


#콘솔에서 미스트 Dapp 띄우기는 명령 :

# ~/go/src/github.com/ethereum/go-ethereum/build/node6/data> open -a /Applications/Ethereum\ Wallet.app --args --rpc ./geth.ipc
