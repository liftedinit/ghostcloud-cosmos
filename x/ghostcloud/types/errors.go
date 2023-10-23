package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/ghostcloud module sentinel errors
var (
	ErrInvalidAddress = sdkerrors.Register(ModuleName, 1100, "invalid address")
)
