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
		DeploymentList: []Deployment{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in deployment
	deploymentIndexMap := make(map[string]struct{})

	for _, elem := range gs.DeploymentList {
		addr, err := sdk.AccAddressFromBech32(elem.Creator)
		if err != nil {
			return err
		}
		index := string(DeploymentKey(addr, elem.Meta.Name))
		if _, ok := deploymentIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for deployment")
		}
		deploymentIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
