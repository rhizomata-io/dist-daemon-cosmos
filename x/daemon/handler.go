package daemon

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon/internal/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgAddMember:
			return handleMsgAddMember(ctx, keeper, msg)
		case MsgRemoveMember:
			return handleMsgRemoveMember(ctx, keeper, msg)
		case MsgHeartbeat:
			return handleMsgHeartbeat(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgAddMember(ctx sdk.Context, keeper Keeper, msg MsgAddMember) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.NodeID)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	member := NewMember (msg.Name, msg.NodeID, msg.Owner )
	keeper.SetMember(ctx, msg.Name, member) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                      // return
}

// Handle a message to delete name
func handleMsgRemoveMember(ctx sdk.Context, keeper Keeper, msg MsgRemoveMember) sdk.Result {
	if !keeper.IsNodePresent(ctx, msg.NodeID) {
		return types.ErrMemberDoesNotExist(types.DefaultCodespace).Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.NodeID)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	
	keeper.RemoveMember(ctx, msg.NodeID)
	return sdk.Result{}
}

// Handle a message to set name
func handleMsgHeartbeat(ctx sdk.Context, keeper Keeper, msg MsgHeartbeat) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.NodeID)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetHeartbeat(ctx, msg.NodeID, msg.Time)
	return sdk.Result{}                      // return
}
