package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/config"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/internal/debug"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"
)

var (
	app = cli.NewApp()

	// Flags
	// keystoreFlag = cli.StringFlag{
	// 	Name:  "keystore",
	// 	Value: filepath.Join(node.DefaultDataDir(), "keystore"),
	// 	Usage: "Directory for the keystore",
	// }
	// configdirFlag = cli.StringFlag{
	// 	Name:  "configdir",
	// 	Value: DefaultConfigDir(),
	// 	Usage: "Directory for Clef configuration",
	// }
	chainIdFlag = cli.Int64Flag{
		Name:  "chainid",
		Value: params.MainnetChainConfig.ChainId.Int64(),
		Usage: "Chain id to use for signing (1=mainnet, 3=Ropsten, 4=Rinkeby, 5=Goerli)",
	}
	listenAddrFlag = cli.StringFlag{
		Name:  "addr",
		Usage: "HTTP server listening interface",
		Value: "localhost:9980",
	}
	// listenPortFlag = cli.IntFlag{
	// 	Name:  "port",
	// 	Usage: "HTTP server listening port",
	// 	Value: 9980,
	// }
	verbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail",
		Value: 3,
	}

	// Commands
	initCommand = cli.Command{
		Action:    utils.MigrateFlags(initializeSecrets),
		Name:      "init",
		Usage:     "Initialize the signer, generate secret storage",
		ArgsUsage: "",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			//configdirFlag,
		},
		Description: `
The init command generates a master seed which Clef can use to store credentials and data needed for
the rule-engine to work.`,
	}
)

func init() {
	app.Name = "Qmanager"
	app.Usage = "Manage Reapchain Qrand operations"
	app.Flags = []cli.Flag{
		//keystoreFlag,
		//configdirFlag,
		chainIdFlag,
		listenAddrFlag,
		// listenPortFlag,
		verbosityFlag,
		utils.DataDirFlag,
		utils.NoUSBFlag,
	}
	app.Commands = []cli.Command{
		initCommand,
	}
	app.Before = func(ctx *cli.Context) error {
		return debug.Setup(ctx)
	}
	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		//prompt.Stdin.Close() // Resets terminal mode.
		return nil
	}
	app.Action = qman
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func qman(ctx *cli.Context) error {
	if args := ctx.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	config := getConfig(ctx)
	path := "./qman"
	if ctx.GlobalIsSet(utils.DataDirFlag.Name) {
		path = ctx.GlobalString(utils.DataDirFlag.Name)
	}
	db, err := ethdb.NewLDBDatabase(filepath.Join(path, "qmandata"), 0, 0)
	if err != nil {
		return errors.New("Cannot open database")
	}
	addr := "127.0.0.1:9980"
	if ctx.GlobalIsSet(listenAddrFlag.Name) {
		addr = ctx.GlobalString(listenAddrFlag.Name)
	}
	qman := NewQmanager(db, config, addr)
	qman.Start()

	return nil
}

func getConfig(ctx *cli.Context) *config.EnvConfig {
	var config config.EnvConfig
	config.GetConfig("REAPCHAIN_ENV", "SETUP_INFO")
	log.Info("Load config.json", "config", config)
	return &config
}

func initializeSecrets(c *cli.Context) error {
	return nil
}
