package main

import (
	"encoding/json"
	// "fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmn "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	
	dapp "github.com/rhizomata-io/dist-daemon-cosmos/daemon"
	// "github.com/rhizomata-io/dist-daemon-cosmos/x/daemon"
	"github.com/cosmos/cosmos-sdk/codec"
	// daemon "github.com/rhizomata-io/dist-daemon-cosmos/daemon"
	
	
	"github.com/rhizomata-io/dist-daemon-cosmos/node"
)

func main() {
	cobra.EnableCommandSorting = false
	
	cdc := dapp.MakeCodec()
	
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()
	
	ctx := server.NewDefaultContext()
	
	rootCmd := &cobra.Command{
		Use:               "daemon",
		Short:             "Daemon App Server",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, dapp.ModuleBasics, dapp.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, dapp.DefaultNodeHome),
		genutilcli.GenTxCmd(
			ctx, cdc, dapp.ModuleBasics, staking.AppModuleBasic{},
			genaccounts.AppModuleBasic{}, dapp.DefaultNodeHome, dapp.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, dapp.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genaccscli.AddGenesisAccountCmd(ctx, cdc, dapp.DefaultNodeHome, dapp.DefaultCLIHome),
		keys.Commands(),
	)
	
	newApp := makeNewAppFn(ctx)
	
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)
	
	rootCmd.AddCommand(
		daemonCmd(ctx, cdc, newApp),
	)
	
	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "NS", dapp.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func daemonCmd( ctx *server.Context, cdc *codec.Codec,  appCreator server.AppCreator)  *cobra.Command  {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Run the full node",
		Long: `Run the full node application with Tendermint in or out of process.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := startDaemon(ctx, cdc, appCreator)
			return err
		},
	}
	cmd.Flags().String("from","operator","name of operator")
	tcmd.AddNodeFlags(cmd)
	return cmd
}

func makeNewAppFn(ctx *server.Context) func(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return func (logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
		// fmt.Println("****** makeNewAppFn :: node.SetNodeKeyFile")
		//
		// cliCtx := context.NewCLIContext().WithCodec(daemon.ModuleCdc)
		//
		// fmt.Println("cliCtx.Client::", cliCtx.Client)
		// fmt.Println("cliCtx.FromAddress::", cliCtx.FromAddress)
		application := dapp.NewDaemonApp(logger, db, baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)))
		application.Init(ctx)
		return application
	}
}

func startDaemon(ctx *server.Context, cdc *codec.Codec, appCreator server.AppCreator) (*tmn.Node, error) {
	cfg := ctx.Config
	home := cfg.RootDir
	traceWriterFile := viper.GetString("trace-store")
	
	db, err := openDB(home)
	if err != nil {
		return nil, err
	}
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return nil, err
	}
	
	application := appCreator(ctx.Logger, db, traceWriter)
	
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}
	
	server.UpgradeOldPrivValFile(cfg)
	
	// create & start tendermint node
	tmNode, err := tmn.NewNode(
		cfg,
		pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(application),
		tmn.DefaultGenesisDocProviderFunc(cfg),
		tmn.DefaultDBProvider,
		tmn.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger.With("module", "node"),
	)
	if err != nil {
		return nil, err
	}
	
	if err := tmNode.Start(); err != nil {
		return nil, err
	}
	
	var cpuProfileCleanup func()
	
	if cpuProfile := viper.GetString("cpu-profile"); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return nil, err
		}
		
		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err := pprof.StartCPUProfile(f); err != nil {
			return nil, err
		}
		
		cpuProfileCleanup = func() {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			f.Close()
		}
	}
	
	node.SetNodeKeyFile(ctx)
	node.SetNode(tmNode)
	node.SetCLIContext(cdc)
	dma := application.(*dapp.App)
	dma.Start()
	
	server.TrapSignal(func() {
		if tmNode.IsRunning() {
			_ = tmNode.Stop()
		}
		
		if cpuProfileCleanup != nil {
			cpuProfileCleanup()
		}
		
		ctx.Logger.Info("exiting...")
	})
	
	// run forever (the node will not be returned)
	select {}
}


func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	
	if height != -1 {
		nsApp := dapp.NewDaemonApp(logger, db)
		err := nsApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	
	nsApp := dapp.NewDaemonApp(logger, db)
	
	return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := sdk.NewLevelDB("application", dataDir)
	return db, err
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile != "" {
		w, err = os.OpenFile(
			traceWriterFile,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0666,
		)
		return
	}
	return
}
