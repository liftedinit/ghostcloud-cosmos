package types_test

import (
	"ghostcloud/testutil/keeper"
	"ghostcloud/testutil/sample"
	"testing"

	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	deployment := sample.CreateDeployment(0, keeper.DATASET_SIZE)
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params:      types.DefaultParams(),
				Deployments: []*types.Deployment{deployment},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicate deployment entry",
			genState: &types.GenesisState{
				Params:      types.DefaultParams(),
				Deployments: []*types.Deployment{deployment, deployment},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
