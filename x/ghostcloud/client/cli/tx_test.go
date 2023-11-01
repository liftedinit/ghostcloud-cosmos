package cli_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"ghostcloud/testutil/network"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/client/cli"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCreateDeployment(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	dir, err := sample.CreateTempDataset()
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	file, err := sample.CreateTempArchive()
	require.NoError(t, err)
	defer file.Close()
	defer os.Remove(file.Name())

	tests := []struct {
		desc   string
		idName string

		fields []string
		args   []string
		err    error
		code   uint32
	}{
		{
			idName: strconv.Itoa(0),

			desc:   "valid dataset",
			fields: []string{dir},
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10))).String()),
				fmt.Sprintf("--%s=%s", cli.FlagDescription, strconv.Itoa(0)),
				fmt.Sprintf("--%s=%s", cli.FlagDomain, strconv.Itoa(0)),
			},
		},
		{
			idName: strconv.Itoa(1),

			desc:   "valid archive",
			fields: []string{file.Name()},
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10))).String()),
				fmt.Sprintf("--%s=%s", cli.FlagDescription, strconv.Itoa(0)),
				fmt.Sprintf("--%s=%s", cli.FlagDomain, strconv.Itoa(0)),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			require.NoError(t, net.WaitForNextBlock())

			args := []string{tc.idName}
			args = append(args, tc.fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateDeployment(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			var resp sdk.TxResponse
			require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, clitestutil.CheckTxCode(net, ctx, resp.TxHash, tc.code))
		})
	}
}
