package ghostcloud

import (
	"math/rand"

	"ghostcloud/testutil/sample"
	ghostcloudsimulation "ghostcloud/x/ghostcloud/simulation"
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = ghostcloudsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgCreateDeployment = "op_weight_msg_deployment"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateDeployment int = 100

	opWeightMsgUpdateDeployment = "op_weight_msg_deployment"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateDeployment int = 100

	opWeightMsgDeleteDeployment = "op_weight_msg_deployment"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteDeployment int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	ghostcloudGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		DeploymentList: []types.Deployment{
			{
				Creator: sample.AccAddress(),
				Meta: &types.DeploymentMeta{
					Name:        "0",
					Description: "0 desc",
					Domain:      "0 domain",
				},
				Files: []*types.File{{
					Meta:    &types.FileMeta{Name: "0 file"},
					Content: []byte{0x00},
				}},
			},
			{
				Creator: sample.AccAddress(),
				Meta: &types.DeploymentMeta{
					Name:        "1",
					Description: "1 desc",
					Domain:      "1 domain",
				},
				Files: []*types.File{{
					Meta:    &types.FileMeta{Name: "1 file"},
					Content: []byte{0x01},
				}},
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&ghostcloudGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateDeployment int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateDeployment, &weightMsgCreateDeployment, nil,
		func(_ *rand.Rand) {
			weightMsgCreateDeployment = defaultWeightMsgCreateDeployment
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateDeployment,
		ghostcloudsimulation.SimulateMsgCreateDeployment(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateDeployment int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateDeployment, &weightMsgUpdateDeployment, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateDeployment = defaultWeightMsgUpdateDeployment
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateDeployment,
		ghostcloudsimulation.SimulateMsgUpdateDeployment(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteDeployment int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteDeployment, &weightMsgDeleteDeployment, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteDeployment = defaultWeightMsgDeleteDeployment
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteDeployment,
		ghostcloudsimulation.SimulateMsgDeleteDeployment(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateDeployment,
			defaultWeightMsgCreateDeployment,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ghostcloudsimulation.SimulateMsgCreateDeployment(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateDeployment,
			defaultWeightMsgUpdateDeployment,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ghostcloudsimulation.SimulateMsgUpdateDeployment(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteDeployment,
			defaultWeightMsgDeleteDeployment,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ghostcloudsimulation.SimulateMsgDeleteDeployment(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}