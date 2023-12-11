package ghostcloud

import (
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, deployment := range genState.Deployments {
		addr := sdk.MustAccAddressFromBech32(deployment.Meta.Creator)
		k.SetDeployment(ctx, addr, deployment.Meta, deployment.Dataset)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.Deployments = keeper.GetAllDeployments(ctx, k)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
