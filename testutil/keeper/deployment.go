package keeper

import (
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateNDeployment(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Deployment {
	items := make([]types.Deployment, n)
	for i := range items {
		items[i].Meta = sample.GetDeploymentMeta(sample.AccAddress(), i)
		items[i].Files = sample.GetDeploymentFiles(i)

		keeper.SetDeployment(ctx, items[i])
	}
	return items
}

func DeploymentFilesToFileMeta(deployment *types.Deployment) []*types.FileMeta {
	meta := make([]*types.FileMeta, len(deployment.Files))
	for i, x := range deployment.Files {
		meta[i] = x.Meta
	}

	return meta
}
