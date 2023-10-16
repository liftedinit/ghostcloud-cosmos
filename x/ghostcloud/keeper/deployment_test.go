package keeper_test

import (
	"strconv"
	"testing"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/testutil/nullify"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNDeployment(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Deployment {
	items := make([]types.Deployment, n)
	for i := range items {
		items[i].Name = strconv.Itoa(i)

		keeper.SetDeployment(ctx, items[i])
	}
	return items
}

func TestDeploymentGet(t *testing.T) {
	keeper, ctx := keepertest.GhostcloudKeeper(t)
	items := createNDeployment(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetDeployment(ctx,
			item.Name,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestDeploymentRemove(t *testing.T) {
	keeper, ctx := keepertest.GhostcloudKeeper(t)
	items := createNDeployment(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveDeployment(ctx,
			item.Name,
		)
		_, found := keeper.GetDeployment(ctx,
			item.Name,
		)
		require.False(t, found)
	}
}

func TestDeploymentGetAll(t *testing.T) {
	keeper, ctx := keepertest.GhostcloudKeeper(t)
	items := createNDeployment(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllDeployment(ctx)),
	)
}
