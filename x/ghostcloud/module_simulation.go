package ghostcloud

import (
	// Using `math/rand` is okay for the simulation, but not in production.
	"math/rand" // #nosec

	keepertest "ghostcloud/testutil/keeper"
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
	// this line is used by starport scaffolding # simapp/module/const
	opWeightMsgCreateDeployment = "op_weight_msg_deployment"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateDeployment int = 100

	opWeightMsgUpdateDeployment = "op_weight_msg_deployment"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateDeployment int = 100

	opWeightMsgRemoveDeployment = "op_weight_msg_deployment"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveDeployment int = 100
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	ghostcloudGenesis := types.GenesisState{
		Params:      types.DefaultParams(),
		Deployments: []*types.Deployment{sample.CreateDeploymentWithAddrAndIndexHtml(accs[0], 0, keepertest.DATASET_SIZE)},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&ghostcloudGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
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

	var weightMsgRemoveDeployment int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRemoveDeployment, &weightMsgRemoveDeployment, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveDeployment = defaultWeightMsgRemoveDeployment
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveDeployment,
		ghostcloudsimulation.SimulateMsgRemoveDeployment(am.accountKeeper, am.bankKeeper, am.keeper),
	))
	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		// this line is used by starport scaffolding # simapp/module/OpMsg
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
			opWeightMsgRemoveDeployment,
			defaultWeightMsgRemoveDeployment,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				ghostcloudsimulation.SimulateMsgRemoveDeployment(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
	}
}
