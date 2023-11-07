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

func TestDeploymentMsgServerUpdate(t *testing.T) {
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

		newMeta := metas[i]
		newMeta.Description = fmt.Sprintf("%s%s", newMeta.Description, "updated")
		newMeta.Domain = fmt.Sprintf("%s%s", newMeta.Domain, "updated")
		_, newPayload := sample.CreateDatasetPayload(100-i, 1)

		_, err = srv.UpdateDeployment(wctx, &types.MsgUpdateDeploymentRequest{Meta: newMeta, Payload: newPayload})
		require.NoError(t, err)
		meta, found := k.GetMeta(ctx, creator, expected.Meta.GetName())
		require.True(t, found)
		require.Equal(t, newMeta, &meta)

		dataset := k.GetDataset(ctx, creator, expected.Meta.GetName())
		require.True(t, found)
		require.Equal(t, newPayload.GetDataset(), dataset)
	}
}
