package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestDeploymentMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateDeployment{Creator: creator,
			Name: strconv.Itoa(i),
		}
		_, err := srv.CreateDeployment(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetDeployment(ctx,
			expected.Name,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestDeploymentMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateDeployment
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateDeployment{Creator: creator,
				Name: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateDeployment{Creator: "B",
				Name: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateDeployment{Creator: creator,
				Name: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.GhostcloudKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateDeployment{Creator: creator,
				Name: strconv.Itoa(0),
			}
			_, err := srv.CreateDeployment(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateDeployment(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetDeployment(ctx,
					expected.Name,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestDeploymentMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteDeployment
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteDeployment{Creator: creator,
				Name: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteDeployment{Creator: "B",
				Name: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteDeployment{Creator: creator,
				Name: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.GhostcloudKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateDeployment(wctx, &types.MsgCreateDeployment{Creator: creator,
				Name: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteDeployment(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetDeployment(ctx,
					tc.request.Name,
				)
				require.False(t, found)
			}
		})
	}
}
