# REAP CHAIN

***[백서 참조](https://reapchain.com/file/en/YellowPaper_ver1.pdf)***

REAP CHAIN intends to improve speed, stability, and security, all at the same time by employing consensus processes which consists of nodes selected by applying quantum random numbers from authorized nodes and general nodes to overcome disadvantages of consensus algorithms previously proposed in various blockchain platforms.



## Contents

- [System Requirements](#System-Requirements)
  - [H/W Specification](#H/W-Specification)
  - [Storage Requirements](#Storage-Requirements)
  - [Operaing System](#Operaing-System)
- [Installation Guide](#Installation-Guide)
  - [Building the source](#Building-the-source)
  - [Executables](#Executables)
- [Startup Endpoint Node](#Startup-Endpoint-Node)
- [Programmatically interfacing geth nodes](#Programmatically-interfacing-geth-nodes)
- [Account Management](#Account-Management)
  - [Managing Accounts](#Managing-Accounts)
  - [Examples](#Examples)
- [License](#License)



---



## System Requirements

For running an REAP CHAIN Endpoint Node requires following specifications.

### H/W Specification

##### Cloud VM

###### Recommended Specification Based on Hanaro-Hosting

| vCPU | Memory(G) | Storage | Network  Bandwidth(Gbps) | Basic  Traffic(G) | Region | VM Total  Price(Won/hour) | VM Total  Price(Won/Month) |
| ---- | --------- | ------- | ------------------------ | ----------------- | ------ | ------------------------- | -------------------------- |
| 4    | 4         | 1,000   | 1                        | 600               | Seoul  | 129                       | 93,120                     |

The information above is from **Hanaro Hosting **and may be changed  by **Hanaro Hosting**.

##### Bare-metal Machine

***[klaytn 내용 참조](https://docs.klaytn.com/node/endpoint-node/system-requirements#bare-metal-machine)***

We do not specify the exact physical machine specification for EN, but any physical machine having hardware configuration similar to the one in the Cloud VM section would be sufficient to operate an EN.



### Storage Requirements

***[klaytn 내용 참조](https://docs.klaytn.com/node/endpoint-node/system-requirements#storage-requirements)***

Assuming 100 TPS in average, 300 bytes average transaction size, and 1-second block latency, the expected EN daily storage requirement is 2.5 GB/day (=300x100x86400).



### Operating System

Recommended environment is Ubuntu 18.04.4 LTS. REAP CHAIN binaries are fully tested on Ubuntu 18.04.4 LTS. macOS binaries are also provided for development purpose.



---



## Installation Guide

### Building the source

##### Linux/Unix

Clone the repository to a directory of your choosing:

```
git clone https://github.com/reapchain/go-reapchain.git
```

Install latest distribution of [Go](https://golang.org/) if you don't have it already.

Building `geth` requires Go and C compilers to be installed:

```
sudo apt-get install -y build-essential
```

Finally, build the `geth` program using the following command.

```
cd go-reapchain
make geth
```

You can now run `build/bin/geth` to start your node.



##### macOS

Clone the repository to a directory of your choosing:

```
git clone https://github.com/reapchain/go-reapchain.git
```

Building `geth` requires the Go compiler:

```
brew install go
```

Finally, build the `geth` program using the following command.

```
cd go-reapchain
make geth
```

If you see some errors related to header files of Mac OS system library, install XCode Command Line Tools, and try again.

```
xcode-select --install
```

You can now run `build/bin/geth` to start your node.



### Executables

| Command    | Description                                                  |
| ---------- | ------------------------------------------------------------ |
| **`geth`** | It is the entry point into the REAP CHAIN network, capable of running as a full node (default), archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the REAP CHAIN network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `geth --help` for command line options. |



---



## Startup Endpoint Node

### Overview

An Endpoint Node has the following roles and functions.

***[참조 링크](https://docs.ethhub.io/using-ethereum/running-an-ethereum-node/)***

##### A full node:

- Stores the full blockchain data available on disk and can serve the network with any data on request.

- Receives new transactions and blocks while participating in block validation.

- Verifies all blocks and states.

- Stores recent state only for more efficient initial sync.

- All state can be derived from a full node.

- Once fully synced, stores all state moving forward.

  

##### A light node:

- Stores the header chain and requests everything else on demand.
- Can verify the validity of the data against the state roots in the block headers.

Light nodes are useful for low capacity devices, such as embedded devices or mobile phones, which can't afford to store multiple dozen Gigabytes of blockchain data.



### Quick Start for running Endpoint Node

Following information is about setting up a REAP CHAIN endpoint node.

##### Endpoint Node directory tree example

```
./node1
├── genesis.json
├── keystore
│   └── node1-account-keystore-file
└── qmanager-nodes.json
```



##### create Endpoint Node directory

```
$ mkdir ./node1
```



##### create account

```
$ geth account new --datadir ./node1

// type passphrase
Your new account is locked with a password. Please give a password. Do not forget this password.
Passphrase:
Repeat passphrase:

// return address
Address: {79180f870c1af86256b024add471f4ed301162d8}
```



##### genesis.json

The genesis block is created using the **genesis state file** or `genesis.json` in Geth. This file contains all the data that will be needed to generate block 0, including who starts out with how much ether.

###### create genesis.json file

```
$ touch ./node1/genesis.json
```

###### REAP CHAIN Main-Network genesis configuration example

Here’s an example of a REAP CHAIN Main-Network genesis state file that initializes block.

```json
{
    "config": {
        "chainId": 2017,
        "homesteadBlock": 1,
        "eip150Block": 2,
        "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "eip155Block": 3,
        "eip158Block": 3,
        "podc": {
            "epoch": 30000,
            "policy": 0
        }
    },
    "nonce": "0x0",
    "timestamp": "0x5e574cb5",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f902bdf90276940585dbbaaeb850ccfec7c90401c24cc31f7a951694090577163d513cdda29106ccc8d71b80d6e0600e940d5d39f8fb9a2139c83af97e1a11749881e97462943305caef9d0195e4ea8f3022a00895543a2d4ca59433c5dfbf34e1b160acddf67375e70f209a65078a9438e62b655c2d767b9f7125d4974fcb5dee93f7c9943accaa82d076d88150365b501d30f996fcd3a49194431d8b17e7e1b766311b70ad336d4851b5b38f2d944c78d3c0441238da9c8ef1548b7e8aeb32b9244a94502775730802f92e7fc44d7ed1da59a5fa1feb0b9456c4e586290b62defeed5c43c8c4e642805fde9e946882edfcd3e95d337966cf1526e338a11a3df43d9477d5bd84fe34b6e716e8d9b8e613e964328442a5947a82f90f9e2e509afe2dd5e66a4781012d4da0e6947bfeaf9a75f1d2af11c3107285e01197ea9f1d239482ecf9e81c2c89751f660e34fc6a0616ba74d86f9484a9986497529a2e217ac021d319c565b16082cf9485b6a2054e660795233979c06601f745b9ac66109499ba3c8d4cc107c5633768c67687eb178aef9081949c7670c1773cd57688efb9bdd92587823df7df2e94a4d4ffaec1965b34360b7f1fdef96883b340cb1f94b007173a6b6c05f3bd3b24d297c6cd3a8423b72594b6761ff2c7e9f5f76d6669eb1b8e38ee9c91806094bf0117d537937e8384f23b5322709d7143436f5894c6119ee3d531a6a67f4c92e85368ec86c333759894cbcc554284c19dd0075b0114b7605732261428e594d11425f73fd6f15cbd9b28b7b90541f904fe27d694e7088f828f59f13f43fa68cbd447c1e99a9b379594e81f66e2a1f34af9cd62890aec815a7724dc031e94fd80e6ead0605d8c4df7b0cabd1c662b036e3292b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0",
    "gasLimit": "0x9886E0",
    "difficulty": "0x1",
    "mixHash": "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
    "coinbase": "0x0000000000000000000000000000000000000000",
    "alloc": {
        "0585dbbaaeb850ccfec7c90401c24cc31f7a9516": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "090577163d513cdda29106ccc8d71b80d6e0600e": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "0d5d39f8fb9a2139c83af97e1a11749881e97462": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "3305caef9d0195e4ea8f3022a00895543a2d4ca5": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "33c5dfbf34e1b160acddf67375e70f209a65078a": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "38e62b655c2d767b9f7125d4974fcb5dee93f7c9": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "3accaa82d076d88150365b501d30f996fcd3a491": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "431d8b17e7e1b766311b70ad336d4851b5b38f2d": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "4c78d3c0441238da9c8ef1548b7e8aeb32b9244a": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "502775730802f92e7fc44d7ed1da59a5fa1feb0b": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "56c4e586290b62defeed5c43c8c4e642805fde9e": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "6882edfcd3e95d337966cf1526e338a11a3df43d": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "77d5bd84fe34b6e716e8d9b8e613e964328442a5": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "7a82f90f9e2e509afe2dd5e66a4781012d4da0e6": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "7bfeaf9a75f1d2af11c3107285e01197ea9f1d23": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "82ecf9e81c2c89751f660e34fc6a0616ba74d86f": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "84a9986497529a2e217ac021d319c565b16082cf": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "85b6a2054e660795233979c06601f745b9ac6610": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "99ba3c8d4cc107c5633768c67687eb178aef9081": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "9c7670c1773cd57688efb9bdd92587823df7df2e": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "a4d4ffaec1965b34360b7f1fdef96883b340cb1f": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "b007173a6b6c05f3bd3b24d297c6cd3a8423b725": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "b6761ff2c7e9f5f76d6669eb1b8e38ee9c918060": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "bf0117d537937e8384f23b5322709d7143436f58": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "c6119ee3d531a6a67f4c92e85368ec86c3337598": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "cbcc554284c19dd0075b0114b7605732261428e5": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "d11425f73fd6f15cbd9b28b7b90541f904fe27d6": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "e7088f828f59f13f43fa68cbd447c1e99a9b3795": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "e81f66e2a1f34af9cd62890aec815a7724dc031e": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "fd80e6ead0605d8c4df7b0cabd1c662b036e3292": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        }
    },
    "number": "0x0",
    "gasUsed": "0x0",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}
```



##### qmanager-nodes.json

***[백서 참조](https://reapchain.com/file/en/YellowPaper_ver1.pdf)***

REAP CHAIN introduces **Qmanager** that generates quantum random numbers that is more secure than pseudorandom numbers in order to implement security enhanced random selection. By using hardware generating quantum random numbers, Qmanager is basically safer than a system that use pseudorandom number generating software algorithms. Qmanager mainly creates and manages quantum random numbers, and manages node information provided by the governance while selecting candidates for steering committee and the coordinator. To make this following qmanager-nodes.json file is needed for REAP CHAIN Main-Network.

###### create qmanager-nodes.json file

```
$ touch ./node1/qmanager-nodes.json
```

###### REAP CHAIN Main-Network Qmanager configuration example

```json
[   	"enode://f09cb0d2104cf26555c39f47654eec233901b2f2066991bf5bc7c5e6ab44165a1a5db73721433c3820f002df45adb42b2d7a44c1ef310011b810df75a2bb9006@116.126.85.24:30501"
]
```



##### init geth with genesis.json

* ###### example

```
$ geth --datadir ./node1 init ./node1/genesis.json
```



##### run geth node

Finally, you can run `geth` node via following example commands. It starts to listening REAP CHAIN Main-Network through `30303` port. Also you can  manage your geth node via JSON-RPC API with `127.0.0.1` IP address and `8545` port.

- ###### example

```
$ geth --datadir ./node1 --networkid 2017 --port 30303 --nat any --rpc --rpcaddr 127.0.0.1 --rpcport 8545 --syncmode full --rpcapi "eth,net,web3,personal,miner,admin,debug" --ipcdisable
```



---



## Programmatically interfacing geth nodes

***[Ethereum README 참조](https://github.com/ethereum/go-reapchain#programmatically-interfacing-geth-nodes)***

As a developer, sooner rather than later you'll want to start interacting with `geth` and the REAP CHAIN network via your own programs and not manually through the console. To aid this, `geth` has built-in support for a JSON-RPC based APIs. These can be exposed via HTTP, WebSockets and IPC (UNIX sockets on UNIX based platforms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by `geth`, whereas the HTTP and WS interfaces need to manually be enabled and only expose a subset of APIs due to security reasons. These can be turned on/off and configured as you'd expect.

HTTP based JSON-RPC API options:

- `--rpc` Enable the HTTP-RPC server
- `--rpcaddr` HTTP-RPC server listening interface (default: `localhost`)
- `--rpcport` HTTP-RPC server listening port (default: `8545`)
- `--rpcapi` API's offered over the HTTP-RPC interface (default: `eth,net,web3`)
- `--rpccorsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
- `--ws` Enable the WS-RPC server
- `--wsaddr` WS-RPC server listening interface (default: `localhost`)
- `--wsport` WS-RPC server listening port (default: `8546`)
- `--wsapi` API's offered over the WS-RPC interface (default: `eth,net,web3`)
- `--wsorigins` Origins from which to accept websockets requests
- `--ipcdisable` Disable the IPC-RPC server
- `--ipcapi` API's offered over the IPC-RPC interface (default: `admin,debug,eth,miner,net,personal,shh,txpool,web3`)
- `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to connect via HTTP, WS or IPC to a `geth` node configured with the above flags and you'll need to speak [JSON-RPC](https://www.jsonrpc.org/specification) on all transports. You can reuse the same connection for multiple requests!



---



## Account Management

***[Ethereum Account Management 참조](https://github.com/ethereum/go-reapchain/wiki/Managing-your-accounts)***

### Managing Accounts

If you lose the password you use to encrypt your account, you will not be able to access that account. Repeat: It is NOT possible to access your account without a password and there is no *forgot my password* option here. Do not forget it.

The CLI `geth` provides account management via the `account` command:

```
$ geth account <command> [options...] [arguments...]
```

Manage accounts lets you create new accounts, list all existing accounts, import a private key into a new account, migrate to newest key format and change your password.

It supports interactive mode, when you are prompted for password as well as non-interactive mode where passwords are supplied via a given password file. Non-interactive mode is only meant for scripted use on test networks or known safe environments.

Make sure you remember the password you gave when creating a new account (with new, update or import). Without it you are not able to unlock your account.

Note that exporting your key in unencrypted format is NOT supported.

Keys are stored under `<DATADIR>/keystore`. Make sure you backup your keys regularly! See [Ethereum DATADIR backup & restore](https://github.com/ethereum/go-reapchain/wiki/Backup-&-restore) for more information. If a custom datadir and keystore option are given the keystore option takes preference over the datadir option.

The newest format of the keyfiles is: `UTC--<created_at UTC ISO8601>-<address hex>`. The order of accounts when listing, is lexicographic, but as a consequence of the timespamp format, it is actually order of creation

It is safe to transfer the entire directory or the individual keys therein between `geth`nodes. Note that in case you are adding keys to your node from a different node, the order of accounts may change. So make sure you do not rely or change the index in your scripts or code snippets.

And again. **DO NOT FORGET YOUR PASSWORD**

```
COMMANDS:
     list    Print summary of existing accounts
     new     Create a new account
     update  Update an existing account
     import  Import a private key into a new account
```

You can get info about subcommands by `geth account <command> --help`.

```
$ geth account list --help
list [command options] [arguments...]

Print a short summary of all accounts

OPTIONS:
  --datadir "/home/bas/.reapchain"  Data directory for the databases and keystore
  --keystore                       Directory for the keystore (default = inside the datadir)
```

Accounts can also be managed via the [Ethereum Javascript Console](https://github.com/ethereum/go-reapchain/wiki/JavaScript-Console)



### Examples

#### Interactive use

##### creating an account

```
$ geth account new
Your new account is locked with a password. Please give a password. Do not forget this password.
Passphrase:
Repeat Passphrase:
Address: {168bc315a2ee09042d83d7c5811b533620531f67}
```

##### Listing accounts in a custom keystore directory

```
$ geth account list --keystore /tmp/mykeystore/
Account #0: {5afdd78bdacb56ab1dad28741ea2a0e47fe41331} keystore:///tmp/mykeystore/UTC--2017-04-28T08-46-27.437847599Z--5afdd78bdacb56ab1dad28741ea2a0e47fe41331
Account #1: {9acb9ff906641a434803efb474c96a837756287f} keystore:///tmp/mykeystore/UTC--2017-04-28T08-46-52.180688336Z--9acb9ff906641a434803efb474c96a837756287f
```

##### Import private key into a node with a custom datadir

```
$ geth account import --datadir /someOtherEthDataDir ./key.prv
The new account will be encrypted with a passphrase.
Please enter a passphrase now.
Passphrase:
Repeat Passphrase:
Address: {7f444580bfef4b9bc7e14eb7fb2a029336b07c9d}
```

##### Account update

```
$ geth account update a94f5374fce5edbc8e2a8697c15331677e6ebf0b
Unlocking account a94f5374fce5edbc8e2a8697c15331677e6ebf0b | Attempt 1/3
Passphrase:
0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b
Account 'a94f5374fce5edbc8e2a8697c15331677e6ebf0b' unlocked.
Please give a new password. Do not forget this password.
Passphrase:
Repeat Passphrase:
0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b
```



#### Non-interactive use

You supply a plaintext password file as argument to the `--password` flag. The data in the file consists of the raw characters of the password, followed by a single newline.

**Note**: Supplying the password directly as part of the command line is not recommended, but you can always use shell trickery to get round this restriction.

```
$ geth account new --password /path/to/password 

$ geth account import  --datadir /someOtherEthDataDir --password /path/to/anotherpassword ./key.prv
```

See [Ethereum Managing Accounts](https://github.com/ethereum/go-reapchain/wiki/Managing-your-accounts) for more information.



---



## License

The go-reapchain library (i.e. all code outside of the `cmd` directory) is licensed under the [GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also included in our repository in the `COPYING.LESSER` file.

The go-reapchain binaries (i.e. all code inside of the `cmd` directory) is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included in our repository in the `COPYING` file.

