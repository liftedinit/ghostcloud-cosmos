package simulation

import (
	"math/rand"

	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgCreateDeploymentArchive(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCreateDeploymentArchive{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the CreateDeploymentArchive simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "CreateDeploymentArchive simulation not implemented"), nil, nil
	}
}
