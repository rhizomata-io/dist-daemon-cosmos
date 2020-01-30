package keeper

import (
	"time"
	
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon/internal/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}


// NewKeeper creates new instances of the daemon Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
	}
}


// Gets the entire Member metadata struct for a name
func (k Keeper) GetMember(ctx sdk.Context, nodeid string) (types.Member, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.IsNodePresent(ctx, nodeid) {
		return types.NilMember, types.ErrMemberDoesNotExist(types.DefaultCodespace)
	}
	bz := store.Get([]byte(nodeid))
	var member types.Member
	k.cdc.MustUnmarshalBinaryBare(bz, &member)
	return member, nil
}


// Check if the node id is present in the store or not
func (k Keeper) IsNodePresent(ctx sdk.Context, nodeid string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(nodeid))
}


// Sets the entire Member metadata struct for a name
func (k Keeper) SetMember(ctx sdk.Context, nodeid string, member types.Member) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(nodeid), k.cdc.MustMarshalBinaryBare(member))
}

// Deletes the entire Member metadata struct for a name
func (k Keeper) RemoveMember(ctx sdk.Context, nodeid string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(nodeid))
}

// Get an iterator over all names in which the keys are the nodeids and the values are the Members
func (k Keeper) GetMembersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}


// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, nodeid string) bool {
	member,err:=k.GetMember(ctx, nodeid)
	if err != nil {
		return false
	}
	return !member.Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, nodeid string) sdk.AccAddress {
	member,err:=k.GetMember(ctx, nodeid)
	if err != nil {
		return nil
	}
	return member.Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, nodeid string, owner sdk.AccAddress) {
	member,err:=k.GetMember(ctx, nodeid)
	if err != nil {
		return
	}
	member.Owner = owner
	k.SetMember(ctx, nodeid, member)
}

// GetHeartbeat - get the current owner of a name
func (k Keeper) GetHeartbeat(ctx sdk.Context, nodeid string) time.Time {
	member,err:=k.GetMember(ctx, nodeid)
	if err != nil {
		return time.Time{}
	}
	return member.Heartbeat
}

// SetHeartbeat - sets the current owner of a name
func (k Keeper) SetHeartbeat(ctx sdk.Context, nodeid string, heartbeat time.Time) {
	member,err:=k.GetMember(ctx, nodeid)
	if err != nil {
		return
	}
	member.Heartbeat = heartbeat
	k.SetMember(ctx, nodeid, member)
}
