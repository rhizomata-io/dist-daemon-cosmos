package daemon

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	MemberRecords []Member `json:"member_records"`
}

func NewGenesisState() GenesisState {
	return GenesisState{MemberRecords:nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.MemberRecords {
		if len(record.NodeID) == 0 {
			return fmt.Errorf("invalid MemberRecord: NodeID: %s. Error: Missing NodeID", record.NodeID)
		}
		if len(record.Name) == 0 {
			return fmt.Errorf("invalid MemberRecord: Name: %s. Error: Missing Name", record.Name)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.MemberRecords {
		keeper.SetMember(ctx, record.NodeID, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Member
	iterator := k.GetMembersIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		nodeid := string(iterator.Key())
		member, err := k.GetMember(ctx, nodeid)
		if err == nil{
			records = append(records, member)
		}
	}
	return GenesisState{MemberRecords: records}
}
