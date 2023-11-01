package types_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"ghostcloud/x/ghostcloud/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	var addr = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()
	// NOTE: We can't use `sample` because it creates a circular dependency.
	deployment := types.Deployment{
		Meta: &types.Meta{
			Creator:     addr,
			Name:        "0",
			Description: "0",
			Domain:      "0",
		},
		Dataset: &types.Dataset{
			Items: []*types.Item{
				{
					Meta:    &types.ItemMeta{Path: "0"},
					Content: &types.ItemContent{Content: []byte("0")},
				},
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
				Params:      types.DefaultParams(),
				Deployments: []*types.Deployment{&deployment},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicate deployment entry",
			genState: &types.GenesisState{
				Params:      types.DefaultParams(),
				Deployments: []*types.Deployment{&deployment, &deployment},
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
