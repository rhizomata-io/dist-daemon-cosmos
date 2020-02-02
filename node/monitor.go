package node

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmn "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/rpc/client"
)

// NodeMonitor
type Monitor struct {
	chainID string
	nodeKeyFile string
	nodeID      string
	node        *tmn.Node
	local      *client.Local
	serverCtx  *server.Context
	cliContext context.CLIContext
	opAddr sdk.AccAddress
}

var monitor = Monitor{}

func SetNodeKeyFile(ctx *server.Context) {
	monitor.serverCtx = ctx
	
	monitor.nodeKeyFile = ctx.Config.NodeKeyFile()
	nodeKey, err := p2p.LoadOrGenNodeKey(monitor.nodeKeyFile)
	if err != nil {
		panic("Cannot load nodeKey File::" + err.Error())
	}
	monitor.nodeID = string(nodeKey.ID())
}

func SetNode(node *tmn.Node) {
	monitor.node = node
	monitor.local = client.NewLocal(node)
	monitor.chainID = node.GenesisDoc().ChainID
}

func GetNodeID() string {
	if len(monitor.nodeID) == 0 {
		panic("NodeID is not set.")
	}
	return monitor.nodeID
}

func GetNode() *tmn.Node {
	return monitor.node
}

func GetLocal() *client.Local {
	return monitor.local
}

func SetCLIContext(codec *codec.Codec) {
	// fmt.Println("FROM::", viper.GetString(flags.FlagFrom))
	
	cliCtx := context.NewCLIContext().WithCodec(codec).WithClient(GetLocal())
	
	cliCtx.BroadcastMode = flags.BroadcastAsync
	// fmt.Println("GetFromAddress::", cliCtx.GetFromAddress())
	
	monitor.cliContext = cliCtx
	
	monitor.opAddr = cliCtx.GetFromAddress()
	
}

func GetCLIContext() context.CLIContext {
	return monitor.cliContext
}

func ChainID() string {
	return monitor.chainID
}

func GetOpAddress() sdk.AccAddress {
	return monitor.opAddr
}