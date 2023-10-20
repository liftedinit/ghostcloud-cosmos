package simulation

import (
	simappparams "cosmossdk.io/simapp/params"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"math/rand"
	"strconv"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func SimulateMsgCreateDeployment(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		i := r.Int()
		msg := &types.MsgCreateDeployment{
			Creator: simAccount.Address.String(),
			Meta:    sample.GetDeploymentMeta(i),
			Files:   sample.GetDeploymentFiles(i),
		}

		_, found := k.GetDeployment(ctx, simAccount.Address, msg.Meta.Name)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Deployment already exist"), nil, nil
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

func SimulateMsgCreateDeploymentArchive(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCreateDeploymentArchive{
			Creator:        simAccount.Address.String(),
			Meta:           sample.GetDeploymentMeta(r.Int()),
			WebsiteArchive: sample.CreateZipWithHTML(),
		}

		_, found := k.GetDeployment(ctx, simAccount.Address, msg.Meta.Name)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Deployment already exist"), nil, nil
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
		var (
			simAccount    = simtypes.Account{}
			deployment    = types.Deployment{}
			msg           = &types.MsgUpdateDeployment{}
			allDeployment = k.GetAllDeployment(ctx)
			found         = false
		)
		for _, obj := range allDeployment {
			simAccount, found = FindAccount(accs, obj.Creator)
			if found {
				deployment = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "deployment creator not found"), nil, nil
		}
		msg.Creator = simAccount.Address.String()
		msg.Meta = sample.GetDeploymentNameMeta(deployment.Meta.Name, r.Int())
		msg.Files = sample.GetDeploymentFiles(r.Int())

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

func SimulateMsgUpdateDeploymentMeta(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount    = simtypes.Account{}
			deployment    = types.Deployment{}
			msg           = &types.MsgUpdateDeploymentMeta{}
			allDeployment = k.GetAllDeployment(ctx)
			found         = false
		)
		for _, obj := range allDeployment {
			simAccount, found = FindAccount(accs, obj.Creator)
			if found {
				deployment = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "deployment creator not found"), nil, nil
		}
		msg.Creator = simAccount.Address.String()
		msg.Meta = sample.GetDeploymentNameMeta(deployment.Meta.Name, r.Int())

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

func SimulateMsgDeleteDeployment(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount    = simtypes.Account{}
			deployment    = types.Deployment{}
			msg           = &types.MsgDeleteDeployment{}
			allDeployment = k.GetAllDeployment(ctx)
			found         = false
		)
		for _, obj := range allDeployment {
			simAccount, found = FindAccount(accs, obj.Creator)
			if found {
				deployment = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "deployment creator not found"), nil, nil
		}
		msg.Creator = simAccount.Address.String()
		msg.Name = deployment.Meta.Name

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
