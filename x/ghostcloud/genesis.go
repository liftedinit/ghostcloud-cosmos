package ghostcloud

import (
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the deployment
	for _, elem := range genState.DeploymentList {
		k.SetDeployment(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.DeploymentList = k.GetAllDeployment(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
