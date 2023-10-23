package keeper_test

import (
	"ghostcloud/testutil/sample"
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
	creator := sample.AccAddress()
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateDeployment{
			Meta:  sample.GetDeploymentMeta(creator, i),
			Files: sample.GetDeploymentFiles(i),
		}
		_, err := srv.CreateDeployment(wctx, expected)
		require.NoError(t, err)
		addr, err := sdk.AccAddressFromBech32(expected.Meta.Creator)
		require.NoError(t, err)
		rst, found := k.GetDeployment(ctx,
			addr,
			expected.Meta.Name,
		)
		require.True(t, found)
		require.Equal(t, expected.Meta, rst.Meta)
	}
}

func TestDeploymentMsgServerCreateArchive(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := sample.AccAddress()
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateDeploymentArchive{
			Meta:           sample.GetDeploymentMeta(creator, i),
			WebsiteArchive: sample.CreateZipWithHTML(),
		}
		_, err := srv.CreateDeploymentArchive(wctx, expected)
		require.NoError(t, err)
		addr, err := sdk.AccAddressFromBech32(expected.Meta.Creator)
		require.NoError(t, err)
		rst, found := k.GetDeployment(ctx,
			addr,
			expected.Meta.Name,
		)
		require.True(t, found)
		require.Equal(t, expected.Meta, rst.Meta)
	}
}

func TestDeploymentMsgServerUpdate(t *testing.T) {
	creator := sample.AccAddress()
	otherCreator := sample.AccAddress()
	tests := []struct {
		desc    string
		request *types.MsgUpdateDeployment
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateDeployment{
				Meta:  sample.GetDeploymentMeta(creator, 0),
				Files: sample.GetDeploymentFiles(0),
			},
		},
		{
			desc: "KeyNotFound - Other Creator",
			request: &types.MsgUpdateDeployment{
				Meta:  sample.GetDeploymentMeta(otherCreator, 0),
				Files: sample.GetDeploymentFiles(0),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "KeyNotFound - Other Meta",
			request: &types.MsgUpdateDeployment{
				Meta:  sample.GetDeploymentMeta(creator, 1000000),
				Files: sample.GetDeploymentFiles(1000000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.GhostcloudKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateDeployment{
				Meta:  sample.GetDeploymentMeta(creator, 0),
				Files: sample.GetDeploymentFiles(0),
			}
			_, err := srv.CreateDeployment(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateDeployment(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				addr, err := sdk.AccAddressFromBech32(expected.Meta.Creator)
				require.NoError(t, err)
				rst, found := k.GetDeployment(ctx,
					addr,
					expected.Meta.Name,
				)
				require.True(t, found)
				require.Equal(t, expected.Meta, rst.Meta)
			}
		})
	}
}

func TestDeploymentMsgServerUpdateMeta(t *testing.T) {
	creator := sample.AccAddress()
	otherCreator := sample.AccAddress()

	tests := []struct {
		desc    string
		request *types.MsgUpdateDeploymentMeta
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateDeploymentMeta{
				Meta: sample.GetDeploymentNameMeta(creator, "0", 1),
			},
		},
		{
			desc: "KeyNotFound - Other Creator",
			request: &types.MsgUpdateDeploymentMeta{
				Meta: sample.GetDeploymentNameMeta(otherCreator, "0", 1),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "KeyNotFound - Other Meta",
			request: &types.MsgUpdateDeploymentMeta{
				Meta: sample.GetDeploymentMeta(creator, 1000000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.GhostcloudKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateDeployment{
				Meta:  sample.GetDeploymentNameMeta(creator, "0", 0),
				Files: sample.GetDeploymentFiles(0),
			}
			_, err := srv.CreateDeployment(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateDeploymentMeta(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				addr, err := sdk.AccAddressFromBech32(expected.Meta.Creator)
				require.NoError(t, err)
				rst, found := k.GetDeployment(ctx,
					addr,
					expected.Meta.Name,
				)
				require.True(t, found)
				require.Equal(t, tc.request.Meta, rst.Meta)
			}
		})
	}
}

func TestDeploymentMsgServerDelete(t *testing.T) {
	creator := sample.AccAddress()
	otherCreator := sample.AccAddress()

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
			desc: "KeyNotFound - Other Creator",
			request: &types.MsgDeleteDeployment{Creator: otherCreator,
				Name: strconv.Itoa(0),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "KeyNotFound - Other Meta",
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

			_, err := srv.CreateDeployment(wctx, &types.MsgCreateDeployment{
				Meta:  sample.GetDeploymentMeta(creator, 0),
				Files: sample.GetDeploymentFiles(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteDeployment(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				addr, err := sdk.AccAddressFromBech32(tc.request.Creator)
				require.NoError(t, err)
				_, found := k.GetDeployment(ctx,
					addr,
					tc.request.Name,
				)
				require.False(t, found)
			}
		})
	}
}
