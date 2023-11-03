package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:      DefaultParams(),
		Deployments: []*Deployment{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	deploymentMetaIndexMap := make(map[string]struct{})
	deploymentFileMetaIndexMap := make(map[string]struct{})

	for _, elem := range gs.Deployments {
		addr, err := sdk.AccAddressFromBech32(elem.Meta.Creator)
		if err != nil {
			return err
		}

		// Check for duplicate meta
		index := string(DeploymentKey(addr, elem.Meta.Name))
		if _, ok := deploymentMetaIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for deployment")
		}
		deploymentMetaIndexMap[index] = struct{}{}

		// Check for duplicate files
		for _, file := range elem.Dataset.Items {
			index = string(DeploymentItemKey(addr, elem.Meta.Name, file.Meta.Path))
			if _, ok := deploymentFileMetaIndexMap[index]; ok {
				return fmt.Errorf("duplicated index for deployment")
			}
			deploymentFileMetaIndexMap[index] = struct{}{}
		}
	}

	return gs.Params.Validate()
}
