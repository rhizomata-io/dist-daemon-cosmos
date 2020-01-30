package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	MemberDoesNotExist sdk.CodeType = 101
)

// ErrMemberDoesNotExist is the error for name not existing
func ErrMemberDoesNotExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, MemberDoesNotExist, "Member does not exist")
}
