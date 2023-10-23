package keeper

import (
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDeployment set a specific deployment in the store from its index
func (k Keeper) SetDeployment(ctx sdk.Context, deployment types.Deployment) {
	addr, _ := sdk.AccAddressFromBech32(deployment.Meta.Creator)
	k.SetDeploymentMeta(ctx, addr, deployment.Meta)
	k.SetDeploymentFiles(ctx, addr, deployment.Meta.Name, deployment.Files)
}

func (k Keeper) SetDeploymentMeta(ctx sdk.Context, addr sdk.AccAddress, meta *types.DeploymentMeta) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	b := k.cdc.MustMarshal(meta)
	store.Set(types.DeploymentKey(addr, meta.Name), b)
}

func (k Keeper) SetDeploymentFiles(ctx sdk.Context, addr sdk.AccAddress, name string, files []*types.File) {
	// NOTE: Safe to ignore the error here because the caller ensures that
	for _, file := range files {
		k.SetDeploymentFileMeta(ctx, addr, name, file)
		k.SetDeploymentFileContent(ctx, addr, name, file)
	}
}

func (k Keeper) SetDeploymentFileMeta(ctx sdk.Context, addr sdk.AccAddress, name string, file *types.File) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentFilesMetaPrefix)
	b := k.cdc.MustMarshal(file.Meta)
	store.Set(types.DeploymentFileKey(addr, name, file.Meta.Name), b)
}

func (k Keeper) SetDeploymentFileContent(ctx sdk.Context, addr sdk.AccAddress, name string, file *types.File) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentFilesContentPrefix)
	b := k.cdc.MustMarshal(file.Content)
	store.Set(types.DeploymentFileKey(addr, name, file.Meta.Name), b)
}

func (k Keeper) GetDeployment(
	ctx sdk.Context,
	addr sdk.AccAddress,
	name string,
) (val types.Deployment, found bool) {
	meta, found := k.GetDeploymentMeta(ctx, addr, name)
	if !found {
		return val, false
	}

	files, found := k.GetDeploymentFiles(ctx, addr, name)
	if !found {
		return val, false
	}

	return types.Deployment{
		Meta:  &meta,
		Files: files}, true
}

func (k Keeper) GetDeploymentView(
	ctx sdk.Context,
	addr sdk.AccAddress,
	name string,
) (val types.DeploymentView, found bool) {
	meta, found := k.GetDeploymentMeta(ctx, addr, name)
	if !found {
		return val, false
	}
	files, found := k.GetDeploymentFileMeta(ctx, addr, name)
	if !found {
		return val, false
	}

	return types.DeploymentView{
		Creator:        addr.String(),
		DeploymentMeta: &meta,
		FilesMeta:      files,
	}, true
}

func (k Keeper) GetDeploymentMeta(
	ctx sdk.Context,
	addr sdk.AccAddress,
	name string,
) (val types.DeploymentMeta, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	b := store.Get(types.DeploymentKey(addr, name))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetDeploymentFiles(
	ctx sdk.Context,
	addr sdk.AccAddress,
	name string,
) (val []*types.File, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentFilesMetaPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.DeploymentKey(addr, name))
	defer iterator.Close()

	var meta types.FileMeta
	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshal(iterator.Value(), &meta)
		content, found := k.GetDeploymentFileContent(ctx, addr, name, meta.Name)
		if !found {
			return val, false
		}
		val = append(val, &types.File{
			Meta:    &meta,
			Content: &content,
		})
	}

	return val, true
}

func (k Keeper) GetDeploymentFileMeta(
	ctx sdk.Context,
	addr sdk.AccAddress,
	name string,
) (val []*types.FileMeta, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentFilesMetaPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.DeploymentKey(addr, name))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var file types.FileMeta
		k.cdc.MustUnmarshal(iterator.Value(), &file)
		val = append(val, &file)
	}

	return val, true
}

func (k Keeper) GetDeploymentFileContent(ctx sdk.Context,
	addr sdk.AccAddress,
	siteName string,
	fileName string) (val types.FileContent, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentFilesContentPrefix)
	b := store.Get(types.DeploymentFileKey(addr, siteName, fileName))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDeployment removes a deployment from the store
func (k Keeper) RemoveDeployment(
	ctx sdk.Context,
	addr sdk.AccAddress,
	name string,

) {
	// Delete deployment meta
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	store.Delete(types.DeploymentKey(addr, name))

	// Delete all file meta
	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentFilesMetaPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.DeploymentKey(addr, name))
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	// And delete all content
	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentFilesContentPrefix)
	iterator = sdk.KVStorePrefixIterator(store, types.DeploymentKey(addr, name))
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// GetAllDeployment returns all deployment.
func (k Keeper) GetAllDeployment(
	ctx sdk.Context,
) (list []types.Deployment) {
	// Let's create an iterator on the deployment meta
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var deploymentMeta types.DeploymentMeta
		k.cdc.MustUnmarshal(iterator.Value(), &deploymentMeta)

		files, found := k.GetDeploymentFiles(ctx, sdk.MustAccAddressFromBech32(deploymentMeta.Creator), deploymentMeta.Name)
		if !found {
			panic("Deployment files not found")
		}

		list = append(list, types.Deployment{
			Meta:  &deploymentMeta,
			Files: files,
		})
	}

	return list
}

func (k Keeper) GetAllDeploymentMeta(
	ctx sdk.Context,
) (list []types.DeploymentMeta) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DeploymentMeta
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
