# Private Network 구성 방법



## Items

#### Binary

- geth: 노드 실행 binary
- qman: Qmanager 실행 binary
- bootnode: 부트노드 실행 binary
- istanbul: genesis.json, 초기 validator 노드들의 nodekey 생성을 위한 tool



#### 설정 파일

- genesis.json: 제네시스 블록 생성을 위한 정보
- config.json: 노드 실행 설정 정보

- static-nodes.json: 부트노드를 사용치 않고 static 연결을 사용할 경우 각 노드의 연결 정보



#### 노드 디렉토리 구조

```
node
├── geth
│   ├── LOCK
│   ├── chaindata/
│   ├── lightchaindata/
│   ├── nodes/
│   ├── nodekey
│   ├── config.json
│   └── static-nodes.json
├── geth.ipc
└── keystore
    └── UTC--2021-03-26T03-38-53.360113400Z--fa99c4ce290b5c8c4c0c063584f89b36e6cc3bdc
```



#### SETUP_INFO

```
setup_info
├── 0/
├── 1/
├── 2/
├── 3/
...
├── config.json
├── genesis.json
├── passwd.txt
└── static-nodes.json

```



## Preparation

#### Prepare binaries for node installation

You can download source code of go-reapchain and istanbul-tool from github.

- go-reapchain: https://github.com/reapchain/go-reapchain
- istanbul-tool: https://github.com/getamis/istanbul-tools

Build go-reapchain and istanbul-tool from source code, then you may get  "geth", "bootnode", "qman" and "istanbul" binaries. Copy them in proper place for excuting them. Otherwise you can use alias.



#### Set SETUP_INFO

Make setup_info directory and specify this as SETUP_INFO in system environment.

```
$ mkdir setup_info		//Directory name could be changed
$ export SETUP_INFO=/.../setup_info
```

> *It can be added in .profile or .bashrc file for permernant use.*



#### Prepare passphrase file (option)

Make passphrase file for the accounts of consensus nodes.

This is just for convenient use during creating and unlock account. It should be used only for test.

```
vi $SETUP_INFO/passwd.txt
-----------------------------------
1234
```



## Create node datas

#### Create keys for bootnode and qmanager

Make keys for bootnode and qmanager and get enode public keys for using when makes config.json.

```
//Create boot.key, qman.key
$ bootnode -genkey $SETUP_INFO/boot.key
$ bootnode -genkey $SETUP_INFO/qman.key	//qman binary can be used

//Get enode public key for enode infos(enode://{PUBKEY}@{IP}:{PORT}) of bootnode and qmanager
$ bootnode -nodekey $SETUP_INFO/boot.key -writeAddress
$ bootnode -nodekey $SETUP_INFO/qman.key -writeAddress
```



#### Make setup datas

Create node keys of consensus nodes, genesis.json, config.json, and static-nodes.json in SETUP_INFO directory.

```
$ cd $SETUP_INFO
$ istanbul setup --num 7 --nodes --verbose --save
$ cd $WORKSPACE
```

> *The "--num" option should be sum of number of senator nodes and candidate nodes.*



#### Create Accounts from node keys

Create accounts from node keys for all nodes.

```
$ geth --datadir $NODE_DIR account import $SETUP_INFO/$i/nodekey --password $SETUP_INFO/passwd.txt
```

> *$i is the directory created by "istanbul setup" command.*



#### Update genesis.json

Change the consensus name under "config" from "istanbul" to "podc" in genesis.json.

```
$ vi $SETUP_INFO/genesis.json
-----------------------------------
{
  "config": {
    "chainId": 7770,
    "homesteadBlock": 1,
    "eip150Block": 2,
    "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "eip155Block": 3,
    "eip158Block": 3,
    "podc": {			//change istanbul to podc
      "epoch": 30000,
      "policy": 0
    }
  },
  "nonce": "0x0",
  "timestamp": "0x605d678c",
  "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f8d9f8939469948bd0243b2ccb523cb3df2220d82a23f0754d9485e5497153ae2f55b56cb71153b40dd70eab08be946fe48948dcb0d19538869c6031b23dfa66eb0ce594e899b6af78b1ac88724ede0867e38336e6bfe19e944372ed27e80db47e7266070d81eabc258c237087946e4d5c343580349f1f82f77b2747250da9551be3943c06cbe144815ff93b68c83ee1037d5202aafd0cb8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0",
  "gasLimit": "0x47b760",
  "difficulty": "0x1",
  "mixHash": "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc": {
    "3c06cbe144815ff93b68c83ee1037d5202aafd0c": {
      "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
    },
    "4372ed27e80db47e7266070d81eabc258c237087": {
      "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
    },
    ...
  },
  "number": "0x0",
  "gasUsed": "0x0",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}
```



#### Make config.json

Add enode informations of bootnodes and qmanagers with enode public key get above.

Split up the accounts that you create accounts above and add into each "senatornodes" and "candidatenodes" field.

The format of config.json is as follow.

```
$ vi $SETUP_INFO/config.json
-----------------------------------
{
  "local": {
    "consensus": {
      "criteria": 4
    },
    "token": {
    },
    "bootnodes": [
      "enode://5d686a07e38d2862322a2b7e829ee......7d0752dec45a7f12c@127.0.0.1:30391"
    ],
    "qmanagers": [
      "enode://9f37019a195a4d490cb09b4acfa67......77c9f73a4cd4cf4c9@127.0.0.1:30500"
    ],
    "senatornodes": [
     "0xca81cd471a007808bc15185357a96e9f09dcb394",
     "0x9a31473d87344330179c6f58c1592ea5bd20e84b",
     "0x39fd383d85dac1850b4b1986d292de51df5ffe8a"
     ],
    "candidatenodes": [
      "baa3e1cba413394a41a1472d945b2ecf841e2773",
      "bb6d81e6d5d68941b78406e68b629b2444b6f923",
      "e84fdd36e9e6635677496f4a2eaeeb099c681ecc",
      "f79f49a4c4068cf0813a3c82bd75d64877283370"
      ]
  },
  "development": {
    ......
  },
  "production": {
    ......
  }
}
```

> *The fields of "bootnodes" and qmanagers" should be filled with enode informations of bootnode and qmanager.*
>
> *The fields of "senatornodes" and "candidatenodes" should be filled with consensus node addresses.*



#### Create node datas with genesis.json

Create node datas with genesis.json for all nodes.

```
$ geth --datadir $NODE_DIR init $SETUP_INFO/genesis.json
```



#### Copy node keys and config.json into each node data directory

```
$ cp -p $SETUP_INFO/$i/nodekey $NODE_DIR/geth/
$ cp -p $SETUP_INFO/config.json $NODE_DIR/geth/
```

> *$i is the directory created by "istanbul setup" command.*



## Run nodes

#### Run qmanager

```
$ qman -qmankey $SETUP_INFO/qman.key -addr $IP:$PORT -verbosity 9
```



#### Run bootnode

```
$ bootnode -nodekey $SETUP_INFO/boot.key -addr $IP:$PORT -verbosity 9
```



#### Run consensus nodes

```
$ geth --networkid 2018 \
	--port $PORT \
	--datadir $NODE_DIR \
	--mine --minerthreads 1 \
	--targetgaslimit 210000000 \
	--unlock 0 --password $SETUP_INFO/passwd.txt \
	--syncmode "full" \
	--rpc \
	--rpcaddr $IP \
	--rpcport $RPCPORT \
	--rpccorsdomain "*" \
	--rpcapi="db,eth,net,web3,personal,miner,admin,debug,ssh,txpool,PoDC" \
	--nat none \
	--verbosity 4
```

> *The "--mine" option should be provided for consensus node.*
>
> *Refer Ethereum command line options for details.* 

