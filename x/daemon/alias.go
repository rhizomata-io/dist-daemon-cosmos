package daemon


import (
	//"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon/internal/keeper"
	"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon/internal/keeper"
	"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
	TStoreKey   = types.TStoreKey
)


var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier           = keeper.NewQuerier
	NewMember           = types.NewMember
	NewMsgAddMember     = types.NewMsgAddMember
	NewMsgRemoveMember  = types.NewMsgRemoveMember
	NewMsgHeartbeat     = types.NewMsgHeartbeat
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
)

type (
	Keeper              = keeper.Keeper
	MsgAddMember        = types.MsgAddMember
	MsgRemoveMember     = types.MsgRemoveMember
	MsgHeartbeat        = types.MsgHeartbeat
	QueryResMember      = types.QueryResMember
	QueryResNodeIDs     = types.QueryResNodeIDs
	Member              = types.Member
)
