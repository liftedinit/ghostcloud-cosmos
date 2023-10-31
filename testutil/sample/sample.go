package sample

import (
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

func CreateNDeployments(n int, datasetSize int) ([]*types.Meta, []*types.Dataset) {
	metas := make([]*types.Meta, n)
	datasets := make([]*types.Dataset, n)
	for i := 0; i < n; i++ {
		metas[i], datasets[i] = CreateDeployment(i, datasetSize)
	}
	return metas, datasets
}

func CreateDeployment(i int, datasetSize int) (*types.Meta, *types.Dataset) {
	return CreateMeta(i), CreateDataset(datasetSize)
}

func CreateMeta(i int) *types.Meta {
	return &types.Meta{
		Creator:     AccAddress(),
		Name:        strconv.Itoa(i),
		Description: strconv.Itoa(i),
		Domain:      strconv.Itoa(i),
	}
}

func CreateDataset(n int) *types.Dataset {
	return &types.Dataset{Items: CreateNItems(n)}
}

func CreateItem(i int) *types.Item {
	return &types.Item{
		Meta:    &types.ItemMeta{Path: strconv.Itoa(i)},
		Content: &types.ItemContent{Content: []byte{byte(i)}},
	}
}

func CreateNItems(n int) []*types.Item {
	items := make([]*types.Item, n)
	for i := 0; i < n; i++ {
		items[i] = CreateItem(i)
	}
	return items
}
