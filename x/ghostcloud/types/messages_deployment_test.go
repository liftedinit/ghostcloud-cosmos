package types

// NOTE: We can't use testutil/sample here because it imports types and we can't have circular imports in tests

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

var addr1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()

func TestMsgCreateDeployment_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateDeployment
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateDeployment{
				Meta: &DeploymentMeta{
					Creator: "invalid_address",
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateDeployment{
				Meta: &DeploymentMeta{
					Creator: addr1,
				},
			},
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

func TestMsgCreateDeploymentArchive_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateDeploymentArchive
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateDeploymentArchive{
				Meta: &DeploymentMeta{
					Creator: "invalid_address",
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateDeploymentArchive{
				Meta: &DeploymentMeta{
					Creator: addr1,
				},
			},
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

func TestMsgUpdateDeployment_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateDeployment
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateDeployment{
				Meta: &DeploymentMeta{
					Creator: "invalid_address",
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateDeployment{
				Meta: &DeploymentMeta{
					Creator: addr1,
				},
			},
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

func TestMsgDeleteDeployment_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteDeployment
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteDeployment{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteDeployment{
				Creator: addr1,
			},
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
