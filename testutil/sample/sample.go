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

func GetDeploymentMeta(i int) *types.DeploymentMeta {
	return &types.DeploymentMeta{
		Name:        strconv.Itoa(i),
		Description: strconv.Itoa(i),
		Domain:      strconv.Itoa(i),
	}
}

func GetDeploymentFiles(i int) []*types.File {
	return []*types.File{
		{
			Meta:    &types.FileMeta{Name: strconv.Itoa(i)},
			Content: make([]byte, i),
		},
	}
}

func GetDeployment(i int) types.Deployment {
	return types.Deployment{
		Creator: AccAddress(),
		Meta:    GetDeploymentMeta(i),
		Files:   GetDeploymentFiles(i),
	}
}

func GetDeploymentList(i int) []types.Deployment {
	deployments := make([]types.Deployment, 0, i)
	for j := 0; j < i; j++ {
		deployments = append(deployments, GetDeployment(j))
	}
	return deployments
}

func GetDuplicatedDeploymentList() []types.Deployment {
	elem := GetDeployment(0)
	return []types.Deployment{elem, elem}
}
