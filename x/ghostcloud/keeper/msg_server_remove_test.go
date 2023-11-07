package keeper_test

import (
	"testing"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestDeploymentMsgServerRemove(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	metas, payloads := sample.CreateNDatasetPayloads(keepertest.NUM_DEPLOYMENT, keepertest.DATASET_SIZE)
	require.Len(t, metas, keepertest.NUM_DEPLOYMENT)

	for i := 0; i < keepertest.NUM_DEPLOYMENT; i++ {
		expected := &types.MsgCreateDeploymentRequest{
			Meta:    metas[i],
			Payload: payloads[i],
		}
		_, err := srv.CreateDeployment(wctx, expected)
		require.NoError(t, err)

		creator, err := sdk.AccAddressFromBech32(expected.Meta.GetCreator())
		require.NoError(t, err)
		retrievedMeta, found := k.GetMeta(ctx, creator, expected.Meta.GetName())
		require.True(t, found)
		require.Equal(t, expected.Meta, &retrievedMeta)

		_, err = srv.RemoveDeployment(wctx, &types.MsgRemoveDeploymentRequest{Creator: expected.Meta.GetCreator(), Name: expected.Meta.GetName()})
		require.NoError(t, err)
		_, found = k.GetMeta(ctx, creator, expected.Meta.GetName())
		require.False(t, found)
	}
}
