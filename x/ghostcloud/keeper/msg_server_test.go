package keeper_test

import (
	"context"
	"ghostcloud/testutil/sample"
	"testing"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func getMetaAndPayload(payloadOption *types.Payload) ([]*types.Meta, []*types.Payload) {
	switch payloadOption.PayloadOption.(type) {
	case *types.Payload_Dataset:
		return sample.CreateNDatasetPayloads(keepertest.NUM_DEPLOYMENT, keepertest.DATASET_SIZE)
	case *types.Payload_Archive:
		return sample.CreateNArchivePayloads(keepertest.NUM_DEPLOYMENT)
	default:
		panic("invalid payload option")
	}
}

func testDeploymentMsgServerCreate(t *testing.T, payloadOption *types.Payload) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	metas, payloads := getMetaAndPayload(payloadOption)
	require.Len(t, metas, keepertest.NUM_DEPLOYMENT)
	require.Len(t, payloads, keepertest.NUM_DEPLOYMENT)

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
	}

}

func TestDeploymentMsgServerCreate_Dataset(t *testing.T) {
	testDeploymentMsgServerCreate(t, &types.Payload{PayloadOption: &types.Payload_Dataset{}})
}

func TestDeploymentMsgServerCreate_Archive(t *testing.T) {
	testDeploymentMsgServerCreate(t, &types.Payload{PayloadOption: &types.Payload_Archive{}})
}
