package simulation

import (
	"archive/zip"
	"bytes"
	"ghostcloud/testutil/sample"
	"io"
	"math/rand"
	"strconv"
	"time"

	simappparams "cosmossdk.io/simapp/params"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Prevent strconv unused error
var _ = strconv.IntSize

const KB = 1024

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

// generateRandomData generates a string with random data.
func generateRandomData(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// CreateZipWithHTML returns a []byte of a zip archive containing an index.html with random data.
func createZipWithHTML() []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Add index.html to zip
	f, err := zipWriter.Create("index.html")
	if err != nil {
		panic(err)
	}

	// Write some random data to index.html
	randomData := generateRandomData(1 + rand.Intn(KB))
	_, err = io.WriteString(f, randomData)
	if err != nil {
		panic(err)
	}

	// Close the zip writer to finish the zip creation
	err = zipWriter.Close()
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
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
			WebsiteArchive: createZipWithHTML(),
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
		msg.Meta = deployment.Meta
		msg.Files = deployment.Files

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
