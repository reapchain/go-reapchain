// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// geth is the official command-line client for Ethereum.
package main

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/reapapis"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/internal/debug"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
)

const (
	clientIdentifier = "chainadapter" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	gitDate   = ""
	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, "Agent from H/L to Reapchain")
	// flags that configure the node
)

func init() {
	// Initialize the CLI app and start Geth
	app.Action = chainadapter
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2013-2020 The go-ethereum Authors"
	app.Commands = []cli.Command{
	}
	app.Flags = append(app.Flags, debug.Flags...)

	app.Before = func(ctx *cli.Context) error {
		return debug.Setup(ctx)
	}
	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// chainadapter is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func chainadapter(ctx *cli.Context) error {
	if args := ctx.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}

	proxyEnvInfos := getProxyEnvInfos()

	// Governance key import
	superAccount := reapapis.NewImportAccount()
	if err := superAccount.Import(proxyEnvInfos.Governance); err != nil {
		log.Error("failed to import super account key", "error", err)
		os.Exit(-1)
	}

	// Reapchain 특정 Node 연결하기
	client, err := rpc.Dial(proxyEnvInfos.Node.RpcAddress)
	if err != nil {
		log.Error("failed to connect with RPC ", "rpc", proxyEnvInfos.Node.RpcAddress, "error", err)
		os.Exit(-1)
	}

	channelInput := make(chan kafka.Message)
	channelOutput := make(chan bool)

	// kafka에 데이터가 들어 왔을 경우, 대신 처리할 콜백 함수
	// result값이 true이면, kafka의 Topic index에 commit 처리함
	callback := func(message kafka.Message) bool {
		log.Debug("receive kafka message", "kafka", string(message.Value))
		channelInput <- message
		result := <-channelOutput
		return result
	}

	// kafka에 연결하기
	myKafka := reapapis.NewKafkaClient([]string{proxyEnvInfos.Kafka.Address}, proxyEnvInfos.Kafka.Topic, callback)

	// kafka 메시지 수신 백그라운드 동작시작
	myKafka.ReadBackground(context.Background(), reapapis.ToEnd)

	for {
		// 데이터가 들어올때 까지 대기
		inputData := <-channelInput
		var dataForProtocol reapapis.ProxyCall
		if err := reapapis.Deserialize(string(inputData.Value), &dataForProtocol); err != nil {
			log.Error("json unmarshalling error, but commit for next Kafka's TX", "err", err)
			channelOutput <- true
		}

		switch dataForProtocol.Call {
		case "burn":
			log.Debug("burn", "protocol", dataForProtocol.Account)
			if exchange(client, superAccount, dataForProtocol.Account, big.NewInt(int64(dataForProtocol.Value)), proxyEnvInfos.Ratio) != nil {
				channelOutput <- false
			}
		default:
			log.Error("unknown Call", "checking", dataForProtocol.Call)
			channelOutput <- false
			continue
		}
		channelOutput <- true
	}
	return nil
}

//REAPCHAIN_ENVFILE, REAPCHAIN_ENVIRON 환경변수로 값 가져오기.
//실패시 프로세스 다운됨
func getProxyEnvInfos() reapapis.ProxyInfo {
	configfile := os.Getenv("REAPCHAIN_ENVFILE")
	if len(configfile) == 0 {
		log.Error("getenv REAPCHAIN_ENVFILE", "empty")
		os.Exit(-1)
	}
	configEnvVar := os.Getenv("REAPCHAIN_ENVIRON")
	if len(configEnvVar) == 0 {
		log.Error("getenv REAPCHAIN_ENVIRON", "empty")
		os.Exit(-1)
	}
	config, err := reapapis.LoadConfigFile(configfile)
	if err != nil {
		log.Error("Get config for reapchain", "error", err)
		os.Exit(-1)
	}

	switch configEnvVar {
	case "local":
		return config.Proxy.Local
	case "test":
		return config.Proxy.Test
	case "production":
		return config.Proxy.Production
	default:
		log.Error("unknown SETUP_INFO", "SETUP_INFO( local|test|production ) must be set", configEnvVar)
	}
	return reapapis.ProxyInfo{}
}

// Reapchain에게 환전요청하는 enpoint API, 내부에서 채번, signing, 환전Tx 송신, Tx의 응답확인함.
// todo gaslimit/gasprice 값은 genesis 로부터 가져올것, eth_gasPrice, eth_estimateGas 함수 통해서 가져올것.
func transfer(client *rpc.Client, account *reapapis.Account, to string, reap *big.Int, ratio float32) error {
	// nonce 채번하기
	var nonce uint64
	var err error
	nonce, err = reapapis.GetNonceTransaction(client, account.Account().Address, nil)
	if err != nil {
		log.Error("failed to get Nonce address : ", account.Account().Address)
		return err
	}

	log.Debug("Nonce value ", "Account", account.Account().Address, "nonce", nonce)

	tx := types.NewTransaction(nonce, account.Account().Address, reap, big.NewInt(3000000), big.NewInt(1), common.FromHex(to))
	tx, err = account.SignTxWithPassphrase(tx, "reapchain", big.NewInt(2017))
	if err != nil {
		log.Error("SignTx error : ", "exchange", err)
		return err
	}

	if err = reapapis.PublishTransaction(client, tx); err != nil {
		log.Error("publishTransaction error : ", "exchange", err)
		return err
	}
	log.Debug("publishTransaction", "Transaction hash", tx.Hash().String())

	var retryCnt uint = 0
	for {
		retryCnt++
		if retryCnt > 10 {
			return errors.New("too many retry")
		}

		time.Sleep(1 * time.Second)

		receipt, err := reapapis.TransactionReceipt(client, tx.Hash())
		if err != nil {
			log.Debug("TransactionReceipt", "retry", retryCnt)
			continue
		} else {
			log.Debug("TransactionReceipt", "Receipt Hash", receipt.String())
			break
		}
	}
	return nil
}

//환전함수
//TODO ratio 환전비율 변수 먹히지 않음
//TODO GasEstimate 30000000 하드코딩되어 있음.
//TODO reapapis.Account --> Wallet으로
//TODO SignTxWithPassphrase, passwd, chainID 2017 하드코딩되어 있음
func exchange(rpc *rpc.Client, account *reapapis.Account, to string, reap *big.Int, ratio float32) error {
	var nonce uint64
	var err error
	nonce, err = getBigestNonce(rpc, account.Account().Address)
	if err != nil {
		log.Error("failed to get Nonce address : ", "Address", account.Account().Address)
		return err
	}
	log.Debug("Nonce value ", "Account", account.Account().Address, "nonce", nonce)

	tx := types.NewTransaction(nonce, account.Account().Address, reap, big.NewInt(3000000), big.NewInt(1), common.FromHex(to))
	tx, err = account.SignTxWithPassphrase(tx, "reapchain", big.NewInt(2017))
	if err != nil {
		log.Error("SignTx error : ", "transfer", err)
		return err
	}

	if err = publishTransaction(rpc, tx); err != nil {
		log.Error("publishTransaction error : ", "transfer", err)
		return err
	}
	log.Debug("publishTransaction", "transfer", tx.Hash().String())

	var retryCnt uint = 0
	for {
		retryCnt++
		if retryCnt > 20 {
			return errors.New("too many retry")
		}

		time.Sleep(1 * time.Second)

		receipt, err := getTransactionRecipt(rpc, tx.Hash())
		if err != nil {
			log.Debug("TransactionReceipt", "retry", retryCnt)
			continue
		} else {
			log.Debug("TransactionReceipt", "Receipt Hash", receipt.TxHash)
			break
		}
	}

	return nil
}

//pending값과, blockchain값중에 큰 nonce값을 가져온다.
func getBigestNonce(rpc *rpc.Client, account common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := ethclient.NewClient(rpc)
	//blockNumber가 nil이면, 최신 블록중에 nonce 값
	nonceFromLastBlock, err := client.NonceAt(ctx, account, nil)
	if err != nil {
		log.Error("get NonceAt", "error", err)
		return 0, err
	}
	nonceFromPending, err := client.PendingNonceAt(ctx, account)
	if err != nil {
		log.Error("Pending Nonce", "error", err)
		return 0, err
	}
	if nonceFromPending > nonceFromLastBlock {
		return nonceFromPending, nil
	}
	return nonceFromLastBlock, nil
}

func publishTransaction(rpc *rpc.Client, tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := ethclient.NewClient(rpc)
	err := client.SendTransaction(ctx, tx)
	if err != nil {
		log.Error("SendTransaction", "error", err)
		return err
	}
	return nil
}

func getTransactionRecipt(rpc *rpc.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var r *types.Receipt
	client := ethclient.NewClient(rpc)

	r, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Error("TransactionReceipt", "error", err)
		return nil, err
	}
	return r, nil
}