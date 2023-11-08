package cli_test

import (
	"os"
	"testing"

	clihelper "ghostcloud/testutil/cli"
	"ghostcloud/testutil/network"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/client/cli"

	"github.com/stretchr/testify/require"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func runRemoveTxTest(t *testing.T, nc *network.Context, tc *network.TxTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		require.NoError(t, nc.Net.WaitForNextBlock())

		args := tc.Args
		out, err := clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdRemoveDeployment(), args)
		if tc.Err == nil {
			require.NoError(t, err)

			var resp sdk.TxResponse
			require.NoError(t, nc.Ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, clitestutil.CheckTxCode(nc.Net, nc.Ctx, resp.TxHash, tc.Code))

			out, err = clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdListDeployments(), []string{})
			require.NoError(t, err)
			require.NotContains(t, out.String(), args[0])
		} else {
			require.Error(t, err)
			require.ErrorContains(t, err, tc.Err.Error())
		}
	})
}

func TestRemoveDeployment(t *testing.T) {
	nc := network.Setup(t)
	commonFlags := network.SetupTxCommonFlags(t, nc)

	testRemove(t, nc, commonFlags)
}

func testRemove(t *testing.T, nc *network.Context, commonFlags []string) {
	newArchive, err := sample.CreateTempArchive(clihelper.IndexHTML, sample.HelloWorldHTMLBody)
	require.NoError(t, err)
	defer os.RemoveAll(newArchive.Name())

	clihelper.CreateDeployment(t, nc, 0, commonFlags)

	runRemoveTxTest(t, nc, &network.TxTestCase{
		Name: "test remove",
		Args: append([]string{"0"}, commonFlags...),
	})
}
