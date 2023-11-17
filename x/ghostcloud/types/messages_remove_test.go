package types_test

import (
	"testing"

	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgRemoveDeployment_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgRemoveDeploymentRequest
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.MsgRemoveDeploymentRequest{Creator: "invalid-addr", Name: "foobar"},
			err:  sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg:  types.MsgRemoveDeploymentRequest{Creator: sample.AccAddress(), Name: "foobar"},
		},
		{
			name: "empty name",
			msg:  types.MsgRemoveDeploymentRequest{Creator: sample.AccAddress(), Name: ""},
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "name has whitespace",
			msg:  types.MsgRemoveDeploymentRequest{Creator: sample.AccAddress(), Name: "foo bar"},
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "name not ascii",
			msg:  types.MsgRemoveDeploymentRequest{Creator: sample.AccAddress(), Name: "foo👍bar"},
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty request",
			msg:  types.MsgRemoveDeploymentRequest{},
			err:  sdkerrors.ErrInvalidAddress,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
