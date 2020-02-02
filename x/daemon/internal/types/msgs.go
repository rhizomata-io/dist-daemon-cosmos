package types

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

type DaemonMsg interface {
	sdk.Msg
	SetOwner(owner sdk.AccAddress)
}

// MsgAddMember add member message
type MsgAddMember struct {
	Name string
	NodeID string
	Owner sdk.AccAddress
}

func NewMsgAddMember(name string, nodeid string) MsgAddMember {
	return MsgAddMember{NodeID:nodeid, Name:name}
}

// Route should return the name of the module
func (msg MsgAddMember) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddMember) Type() string { return "add_member" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddMember) ValidateBasic() sdk.Error {
	if len(msg.Name) == 0 || len(msg.NodeID) == 0 {
		return sdk.ErrUnknownRequest("Name and/or ID cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddMember) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddMember) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgAddMember) SetOwner(owner sdk.AccAddress) {
	msg.Owner = owner
}


// MsgRemoveMember remove member message
type MsgRemoveMember struct {
	NodeID string
	Owner sdk.AccAddress
}

func NewMsgRemoveMember(nodeid string) MsgRemoveMember {
	return MsgRemoveMember{NodeID:nodeid}
}

// Route should return the name of the module
func (msg MsgRemoveMember) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRemoveMember) Type() string { return "remove_member" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRemoveMember) ValidateBasic() sdk.Error {
	if len(msg.NodeID) == 0 {
		return sdk.ErrUnknownRequest("ID cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRemoveMember) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRemoveMember) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgRemoveMember) SetOwner(owner sdk.AccAddress) {
	msg.Owner = owner
}

// MsgHeartbeat heartbeat message
type MsgHeartbeat struct {
	NodeID      string
	Time    time.Time
	Owner sdk.AccAddress
}

func NewMsgHeartbeat(nodeid string, time time.Time, owner sdk.AccAddress) MsgHeartbeat {
	return MsgHeartbeat{NodeID: nodeid, Time: time, Owner:owner}
}

// Route should return the name of the module
func (msg MsgHeartbeat) Route() string { return RouterKey }

// Type should return the action
func (msg MsgHeartbeat) Type() string { return "heartbeat" }

// ValidateBasic runs stateless checks on the message
func (msg MsgHeartbeat) ValidateBasic() sdk.Error {
	if len(msg.NodeID) == 0 {
		return sdk.ErrUnknownRequest("ID cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgHeartbeat) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgHeartbeat) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgHeartbeat) SetOwner(owner sdk.AccAddress) {
	msg.Owner = owner
}
