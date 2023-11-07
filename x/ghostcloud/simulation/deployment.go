package simulation

import (
	"math/rand"
	"strconv"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"

	simappparams "cosmossdk.io/simapp/params"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgCreateDeployment(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		i := r.Intn(5)
		meta, payload := sample.CreateDatasetPayloadWithAddr(simAccount.Address.String(), i, keepertest.DATASET_SIZE)
		msg := &types.MsgCreateDeploymentRequest{
			Meta:    meta,
			Payload: payload,
		}

		found := k.HasDeployment(ctx, simAccount.Address, meta.GetName())
		if found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "deployment already exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgUpdateDeployment(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		i := r.Intn(5)
		meta, payload := sample.CreateDatasetPayloadWithAddr(simAccount.Address.String(), i, keepertest.DATASET_SIZE)
		msg := &types.MsgUpdateDeploymentRequest{
			Meta:    meta,
			Payload: payload,
		}

		found := k.HasDeployment(ctx, simAccount.Address, meta.GetName())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "deployment doesn't exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgRemoveDeployment(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		i := r.Intn(5)
		msg := &types.MsgRemoveDeploymentRequest{
			Creator: simAccount.Address.String(),
			Name:    strconv.Itoa(i),
		}

		found := k.HasDeployment(ctx, simAccount.Address, strconv.Itoa(i))
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "deployment doesn't exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
