package reapapis

// 이 api는 공통적으로 rpc.Client를 인자를 필요로 한다.
// Agent가 이 api를 사용할 때, 고민할 부분
// 1. Reapchain의 분산원장의 동기화부분에서 시간차가 발생된다.
// 2. Tx는 한 노드에 최초로 수신되나, 이 노드에서 다른 노드에서 전파되는 시간이 필요하다.
// 3. 이 때. Agent가 항상 바로보고 있는 Node가 다운 될 경우,
// 4. 다른 노드에서 해당되는 Tx의 값이 항상 과거의 보냈던 최초 Tx가 수신되었다고 보장할 수 없다.
// 시나리오
// Agent ----> A계정 Tx1 ----> Node1(Tx1 수신완료) ---전파--> other NodeX
// Node1 사라짐.
// Agent ----> A계정 Tx2 ----> another Node2 (nonce 채번.. 등)
// another Node2는 자기의 원장(pending) 기준으로 nonce 채번
// Agent ----> Tx2`(환전요청) ----> another Node2
// Agent nonce 실패 발생됨 이때 에러는 어떻게 처리해야 하는가?
// Kafka에서 Tx은 순서대로 들어온다는 보장이 없다. ????
// Agent간끼리 IPC 필요한가?
// 모든게 정상적이나, Block의 생성이 10초이상 지연될 경우
// 수백명의 환전요청할 경우(선형처리라 문제가 있음)

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"
)

func TransactionReceipt(client *rpc.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var r *types.Receipt
	err := client.CallContext(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	return r, err
}

// todo 이 함수도 ethclient.go에 비슷한 기능이 함수가 있음
func PublishTransaction(client *rpc.Client, tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	return client.CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(data))
}

// blockchain, pending area에서 nonce 값을 비교하여 큰 값을 가져 온다.
// todo 이 함수도 ethclient.go에 비슷한 기능이 함수가 있음
func GetNonceTransaction(client *rpc.Client, account common.Address, blockNumber *big.Int) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var resultFromBlock, resultFromPending hexutil.Uint64
	err := client.CallContext(ctx, &resultFromBlock, "eth_getTransactionCount", account, toBlockNumArg(blockNumber))
	if err != nil {
		log.Error("eth_getTransactionCount to Blockchain", "err", account.String())
		return 0, nil
	}

	err = client.CallContext(ctx, &resultFromPending, "eth_getTransactionCount", account, "pending")
	if err != nil {
		log.Error("eth_getTransactionCount to pending area", "err", account.String())
		return 0, nil
	}
	if resultFromPending > resultFromBlock {
		return uint64(resultFromPending), nil
	}
	return uint64(resultFromBlock), nil
}

// todo 이 함수는 ethclient.go에 중복으로 선언되어 있음. 정리가 필요.
func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

type rpcApi struct {
	conn *rpc.Client
}

//ex) NewHttpRpc("http://192.168.0.15:8541")
func NewHttpRpc(endpoint string) (*rpcApi, error) {
	httpRpcClient, err := rpc.Dial(endpoint)
	//make(chan *types.Header, 16)
	if err != nil {
		log.Error("failed to Connect", "rpc.Dail error", err)
		return nil, err
	}
	return &rpcApi{httpRpcClient}, nil
}

//
//	type Block struct {
//		Number     *hexutil.Big
//		ParentHash *hexutil.Big
//	}
//  var result Block
// ex) CallRpcIntoData(&Block, "eth_getBlockByNumber", "latest", false)
func (r *rpcApi) callRpc(rpcFunc string, args ...interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var resultType interface{}
	err := r.conn.CallContext(ctx, &resultType, rpcFunc, args)
	if err != nil {
		log.Error("test subscribeBlocks()", "cant get latest block", err)
		return resultType, err
	}
	return resultType, err
}

//TODO: 처리하다가 data가 문제가 있거나, 혹은 예측하지 하지 않는 데이터 입수시, processingFunc 실패가 떨어진다.
//처리하다가 실패했을 경우, 향후 처리가능하도록 개선해야 됨. 어쩔수 없이 결과 모니터링이 필요
//이 함수 진입하게 되면, hold 된다(백그라운드에서 계속적인 처리).
func (r *rpcApi) GoRpcToProcessing(processingFunc func(data interface{}) bool, rpcFunc string, args ...interface{}) {
	subChan := make(chan interface{})

	go func() {
		for i := 0; ; i++ {
			if i > 0 {
				time.Sleep(2 * time.Second)
			}
			someTypeResult, err := r.callRpc(rpcFunc, args)
			if err != nil {
				log.Error("Error Rpc call", rpcFunc, err)
			}
			subChan <- someTypeResult
		}
	}()

	// Print events from the subscription as they arrive.
	for data := range subChan {
		if !processingFunc(data) {
			log.Error("failed to process Data", data, false)
		}
	}
}
