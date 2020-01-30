package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	
	"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon/internal/types"
)

// query endpoints supported by the nameservice Querier
const (
	QueryMember   = "member"
	QueryNodeIDs   = "nodeids"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryMember:
			return queryMember(ctx, path[1:], req, keeper)
		case QueryNodeIDs:
			return queryNodeIDs(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}


// nolint: unparam
func queryMember(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	member, err := keeper.GetMember(ctx, path[0])
	
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	
	res, err := codec.MarshalJSONIndent(keeper.cdc, member)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	
	return res, nil
}

func queryNodeIDs(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var namesList types.QueryResNodeIDs
	
	iterator := keeper.GetMembersIterator(ctx)
	
	for ; iterator.Valid(); iterator.Next() {
		namesList = append(namesList, string(iterator.Key()))
	}
	
	res, err := codec.MarshalJSONIndent(keeper.cdc, namesList)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	
	return res, nil
}

