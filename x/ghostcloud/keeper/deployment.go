package keeper

import (
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDeployment set a specific deployment in the store from its index
func (k Keeper) SetDeployment(ctx sdk.Context, deployment types.Deployment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentKeyPrefix)

	// NOTE: Safe to ignore the error here because the caller ensures that
	addr, _ := sdk.AccAddressFromBech32(deployment.Creator)
	name := deployment.Meta.Name
	b := k.cdc.MustMarshal(&deployment)
	store.Set(types.DeploymentKey(addr, name), b)
}

// GetDeployment returns a deployment from its index
func (k Keeper) GetDeployment(
	ctx sdk.Context,
	addr sdk.AccAddress,
	name string,
) (val types.Deployment, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentKeyPrefix)
	b := store.Get(types.DeploymentKey(addr, name))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetDeploymentFile(ctx sdk.Context,
	addr sdk.AccAddress,
	siteName string,
	fileName string) (val *types.File, found bool) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentKeyPrefix)
	b := store.Get(types.DeploymentKey(addr, siteName))
	if b == nil {
		return val, false
	}

	var deployment types.Deployment
	k.cdc.MustUnmarshal(b, &deployment)

	for _, file := range deployment.Files {
		if file.Name == fileName {
			return file, true
		}
	}

	return nil, false
}

// RemoveDeployment removes a deployment from the store
func (k Keeper) RemoveDeployment(
	ctx sdk.Context,
	addr sdk.AccAddress,
	name string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentKeyPrefix)
	store.Delete(types.DeploymentKey(addr, name))
}

// GetAllDeployment returns all deployment
func (k Keeper) GetAllDeployment(ctx sdk.Context) (list []types.Deployment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Deployment
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
