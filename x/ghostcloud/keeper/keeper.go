package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"ghostcloud/x/ghostcloud/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) HasDeployment(ctx sdk.Context, creator sdk.AccAddress, name string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	return store.Has(types.DeploymentKey(creator, name))
}

func (k Keeper) SetDeployment(ctx sdk.Context, meta *types.Meta, dataset *types.Dataset) {
	addr, _ := sdk.AccAddressFromBech32(meta.GetCreator())
	k.SetMeta(ctx, addr, meta)
	k.SetDataset(ctx, addr, meta.GetName(), dataset)
}

func (k Keeper) SetMeta(ctx sdk.Context, addr sdk.AccAddress, meta *types.Meta) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	b := k.cdc.MustMarshal(meta)
	store.Set(types.DeploymentKey(addr, meta.GetName()), b)
}

func (k Keeper) SetDataset(ctx sdk.Context, addr sdk.AccAddress, name string, dataset *types.Dataset) {
	// NOTE: Safe to ignore the error here because the caller ensures that
	for _, item := range dataset.GetItems() {
		k.SetItem(ctx, addr, name, item)
	}
}

func (k Keeper) Remove(ctx sdk.Context, addr sdk.AccAddress, name string) {
	k.RemoveDataset(ctx, addr, name)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	store.Delete(types.DeploymentKey(addr, name))
}

func (k Keeper) RemoveDataset(ctx sdk.Context, addr sdk.AccAddress, name string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentItemMetaPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.DeploymentKey(addr, name))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentItemContentPrefix)
	iterator = sdk.KVStorePrefixIterator(store, types.DeploymentKey(addr, name))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) SetItem(ctx sdk.Context, addr sdk.AccAddress, name string, item *types.Item) {
	// Set Item meta
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentItemMetaPrefix)

	meta := item.GetMeta()
	path := meta.GetPath()

	b := k.cdc.MustMarshal(meta)
	store.Set(types.DeploymentItemKey(addr, name, path), b)

	// Set Item content
	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentItemContentPrefix)
	b = k.cdc.MustMarshal(item.GetContent())
	store.Set(types.DeploymentItemKey(addr, name, path), b)
}

func (k Keeper) GetDataset(ctx sdk.Context, addr sdk.AccAddress, name string) (dataset *types.Dataset) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentItemMetaPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.DeploymentKey(addr, name))
	defer iterator.Close()

	items := make([]*types.Item, 0)
	for ; iterator.Valid(); iterator.Next() {
		var meta types.ItemMeta
		k.cdc.MustUnmarshal(iterator.Value(), &meta)

		store = prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentItemContentPrefix)
		b := store.Get(types.DeploymentItemKey(addr, name, meta.GetPath()))
		if b == nil {
			continue
		}

		var content types.ItemContent
		k.cdc.MustUnmarshal(b, &content)

		items = append(items, &types.Item{
			Meta:    &meta,
			Content: &content,
		})
	}

	return &types.Dataset{
		Items: items,
	}
}

func (k Keeper) GetMeta(ctx sdk.Context, addr sdk.AccAddress, name string) (meta types.Meta, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentMetaKeyPrefix)
	b := store.Get(types.DeploymentKey(addr, name))
	if b == nil {
		return meta, false
	}

	k.cdc.MustUnmarshal(b, &meta)
	return meta, true
}

func (k Keeper) GetItemContent(ctx sdk.Context, addr sdk.AccAddress, name string, path string) (content types.ItemContent, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DeploymentItemContentPrefix)
	b := store.Get(types.DeploymentItemKey(addr, name, path))
	if b == nil {
		return content, false
	}

	k.cdc.MustUnmarshal(b, &content)
	return content, true
}

func (k Keeper) GetAllMeta(ctx sdk.Context) (metas []*types.Meta) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DeploymentMetaKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var meta types.Meta
		k.cdc.MustUnmarshal(iterator.Value(), &meta)

		metas = append(metas, &meta)
	}

	return metas
}
