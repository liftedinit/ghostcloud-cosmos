package keeper_test

import (
	"fmt"
	"testing"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func testDeploymentMsgServerRemove(t *testing.T, k *keeper.Keeper, ctx sdk.Context, tc keepertest.MsgServerTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		srv := keeper.NewMsgServerImpl(*k)
		wctx := sdk.WrapSDKContext(ctx)

		for _, meta := range tc.Metas {
			_, err := srv.RemoveDeployment(wctx, &types.MsgRemoveDeploymentRequest{Creator: meta.GetCreator(), Name: meta.GetName()})
			if tc.Err == nil {
				require.NoError(t, err)
				creator, err := sdk.AccAddressFromBech32(meta.GetCreator())
				require.NoError(t, err)
				_, found := k.GetMeta(ctx, creator, meta.GetName())
				require.False(t, found)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.Err.Error())
			}
		}
	})
}

func testDeploymentMsgServerRemoveValid(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	tc := keepertest.MsgServerTestCase{
		Name:  "remove_valid",
		Metas: metas,
	}
	testDeploymentMsgServerRemove(t, k, ctx, tc)
}

func testDeploymentMsgServerRemoveEmptyName(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta := sample.CreateMeta(0)
	meta.Name = ""
	tc := keepertest.MsgServerTestCase{
		Name:  "remove_empty_name",
		Metas: []*types.Meta{meta},
		Err:   fmt.Errorf(types.NameShouldNotBeEmpty),
	}
	testDeploymentMsgServerRemove(t, k, ctx, tc)
}

func testDeploymentMsgServerRemoveEmptyCreator(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta := sample.CreateMeta(0)
	meta.Creator = ""
	tc := keepertest.MsgServerTestCase{
		Name:  "remove_invalid_creator",
		Metas: []*types.Meta{meta},
		Err:   fmt.Errorf(types.CreatorShouldNotBeEmpty),
	}
	testDeploymentMsgServerRemove(t, k, ctx, tc)
}

func testDeploymentMsgServerRemoveInvalidCreator(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta := sample.CreateMeta(0)
	meta.Creator = "invalid creator"
	tc := keepertest.MsgServerTestCase{
		Name:  "remove_invalid_creator",
		Metas: []*types.Meta{meta},
		Err:   fmt.Errorf("invalid creator address"),
	}
	testDeploymentMsgServerRemove(t, k, ctx, tc)
}

func testDeploymentMsgServerRemoveNonExisting(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta := sample.CreateMeta(0)
	tc := keepertest.MsgServerTestCase{
		Name:  "remove_non_existing",
		Metas: []*types.Meta{meta},
		Err:   fmt.Errorf("unable to remove this deployment"),
	}
	testDeploymentMsgServerRemove(t, k, ctx, tc)
}

func TestDeploymentMsgServerRemove(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)

	testDeploymentMsgServerRemoveValid(t, k, ctx)
	testDeploymentMsgServerRemoveEmptyName(t, k, ctx)
	testDeploymentMsgServerRemoveEmptyCreator(t, k, ctx)
	testDeploymentMsgServerRemoveInvalidCreator(t, k, ctx)
	testDeploymentMsgServerRemoveNonExisting(t, k, ctx)
}
