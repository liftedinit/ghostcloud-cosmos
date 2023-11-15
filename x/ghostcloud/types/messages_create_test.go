package types_test

import (
	"strings"
	"testing"

	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stretchr/testify/require"
)

func TestMsgCreateDeployment_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgCreateDeploymentRequest
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.MsgCreateDeploymentRequest{Meta: sample.CreateMetaInvalidAddress()},
			err:  sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg:  types.MsgCreateDeploymentRequest{Meta: sample.CreateMeta(0)},
		},
		{
			name: "empty name",
			msg:  types.MsgCreateDeploymentRequest{Meta: &types.Meta{Creator: sample.AccAddress(), Name: ""}},
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "name with whitespace",
			msg:  types.MsgCreateDeploymentRequest{Meta: &types.Meta{Creator: sample.AccAddress(), Name: "name with whitespace"}},
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "name with non-ascii",
			msg:  types.MsgCreateDeploymentRequest{Meta: &types.Meta{Creator: sample.AccAddress(), Name: "你好"}},
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid domain",
			msg:  types.MsgCreateDeploymentRequest{Meta: &types.Meta{Creator: sample.AccAddress(), Name: "name", Domain: "invalid domain"}},
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "domain too long",
			msg:  types.MsgCreateDeploymentRequest{Meta: &types.Meta{Creator: sample.AccAddress(), Name: "name", Domain: strings.Repeat("a", 65)}},
			err:  sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty request",
			msg:  types.MsgCreateDeploymentRequest{},
			err:  sdkerrors.ErrInvalidRequest,
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
