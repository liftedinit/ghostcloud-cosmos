package types_test

import (
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"testing"
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
