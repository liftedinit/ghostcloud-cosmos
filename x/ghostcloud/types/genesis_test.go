package types_test

// NOTE: We can't use testutil/sample here because it imports types and we can't have circular imports in tests

import (
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
)

var addr1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()

func TestGenesisState_Validate(t *testing.T) {
	deployment := types.Deployment{
		Creator: addr1,
		Meta: &types.DeploymentMeta{
			Name:        "name",
			Description: "description",
			Domain:      "domain",
		},
		Files: []*types.File{
			{
				Meta:    &types.FileMeta{Name: "name"},
				Content: []byte{1},
			},
		},
	}
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

				DeploymentList: []types.Deployment{
					deployment,
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated deployment",
			genState: &types.GenesisState{
				DeploymentList: []types.Deployment{
					deployment, deployment,
				},
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
