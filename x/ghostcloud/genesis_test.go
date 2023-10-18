package ghostcloud_test

import (
	"ghostcloud/testutil/sample"
	"testing"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/testutil/nullify"
	"ghostcloud/x/ghostcloud"
	"ghostcloud/x/ghostcloud/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		DeploymentList: []types.Deployment{
			{
				Creator: sample.AccAddress(),
				Meta:    sample.GetDeploymentMeta(0),
				Files:   sample.GetDeploymentFiles(0),
			},
			{
				Creator: sample.AccAddress(),
				Meta:    sample.GetDeploymentMeta(1),
				Files:   sample.GetDeploymentFiles(1),
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.GhostcloudKeeper(t)
	ghostcloud.InitGenesis(ctx, *k, genesisState)
	got := ghostcloud.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.DeploymentList, got.DeploymentList)
	// this line is used by starport scaffolding # genesis/test/assert
}
