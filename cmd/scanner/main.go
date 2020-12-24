package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/internal/debug"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"
	"math/big"
	"os"
	"time"
)

// Commonly used command line flags.
var (
	gitCommit = ""
	gitDate   = ""
	app       = utils.NewApp(gitCommit, "the Scanner for Reapchain")

//mysqlConnect = cli.StringFlag{
//	Name:  "passwordfile",
//	Usage: "the file that contains the password for the keyfile",
//}
//jsonFlag = cli.BoolFlag{
//	Name:  "json",
//	Usage: "output JSON instead of human-readable format",
//}
)

func init() {
	app.Action = scanner
	app.Commands = []cli.Command{}

	// log 기능 추가
	app.Flags = append(app.Flags, debug.Flags...)

	app.Before = func(ctx *cli.Context) error {
		return debug.Setup(ctx)
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}

	log.Info("Go Scanner.")
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func scanner(ctx *cli.Context) error {
	var err error

	rpcDalier, err := rpc.Dial("http://192.168.11.1:8541")
	if err != nil {
		log.Error("rpc", "rpc.Dial error", err)
		os.Exit(-1)
	}

	client := ethclient.NewClient(rpcDalier)
	block, err := client.BlockByNumber(context.Background(),big.NewInt(1))
	if err != nil {
		log.Error("rcp", "eth_getBlockByNumber", err)
		os.Exit(-1)
	}

	log.Debug("Block", "Block Content", block.Header().ParentHash)

	return nil
}

func callbackfunc(result interface{}) bool {
	value, ok := result.(string)
	if !ok {
		log.Error("interface casting error", "ret=", ok)
		os.Exit(-1)
	}
	v, _ := hexutil.DecodeBig(value)
	log.Debug("Result", "value", v)
	return true
}

type Block struct {
	blocknumber     uint64    //| bigint unsigned | NO   | PRI | NULL    |       |
	timestamp       time.Time //| timestamp       | NO   |     | NULL    |       |
	miner           string    //| char(42)        | NO   |     | NULL    |       |
	blockreward     float64   //| double(40,20)   | YES  |     | NULL    |       |
	unclesreward    float64   //| double(40,20)   | YES  |     | NULL    |       |
	difficulty      uint64    //| bigint unsigned | YES  |     | NULL    |       |
	totaldifficulty uint64    //| bigint unsigned | YES  |     | NULL    |       |
	size            uint64    //| bigint unsigned | NO   |     | NULL    |       |
	gasused         uint64    //| bigint unsigned | NO   |     | NULL    |       |
	gaslimit        uint64    //| bigint unsigned | NO   |     | NULL    |       |
	extradata       string    //| mediumtext      | YES  |     | NULL    |       |
	hash            string    //| char(66)        | NO   |     | NULL    |       |
	parenthash      string    //| char(66)        | NO   |     | NULL    |       |
	sha3uncles      string    //| char(66)        | YES  |     | NULL    |       |
	stateroot       string    //| char(66)        | YES  |     | NULL    |       |
	nonce           string    //| char(18)        | YES  |     | NULL    |       |
}

// Block struct의 타입의 채널에 데이터가 수신되면,
// Block table에 데이터를 저장한다.
func insertBlockTable(data <-chan Block) {
	db, err := sql.Open("mysql", "reapscanner:45b8d01caf660bf63e9b69fb13ab6e40@tcp(192.168.0.97:3306)/reapdb")
	if err != nil {
		log.Error("Connect Mysql", "sql.Open error", err)
	}
	defer db.Close()

	for {
		blockData := <-data
		result, err := db.Exec("insert into block(blocknumber, timestamp, miner, size, gasused, gaslimit, hash, parenthash)",
			blockData.blocknumber,
			blockData.timestamp,
			blockData.miner,
			blockData.size,
			blockData.gasused,
			blockData.gaslimit,
			blockData.hash,
			blockData.parenthash)

		if err != nil {
			log.Error("insert query", "error", err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			log.Error("check result after insert", "error", err)
		}
		if rows != 1 {
			log.Error("the number of rows is too many", "error", "expected to affect 1 row, affected", "row", rows)
		}
	}
}