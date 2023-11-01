package keeper

import (
	"ghostcloud/x/ghostcloud/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetAllDeployments(ctx sdk.Context, k Keeper) (deployments []*types.Deployment) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DeploymentMetaKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var meta types.Meta
		k.cdc.MustUnmarshal(iterator.Value(), &meta)

		creator := sdk.MustAccAddressFromBech32(meta.GetCreator())
		dataset := k.GetDataset(ctx, creator, meta.GetName())

		deployments = append(deployments, &types.Deployment{
			Meta:    &meta,
			Dataset: dataset,
		})
	}

	return
}
