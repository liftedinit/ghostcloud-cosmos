package cli_test

import (
	"fmt"
	"os"
	"testing"

	"ghostcloud/testutil/network"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/client/cli"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func runUpdateTxTest(t *testing.T, nc *network.Context, tc *network.TxTestCase, expected *types.Deployment) {
	t.Run(tc.Name, func(t *testing.T) {
		require.NoError(t, nc.Net.WaitForNextBlock())

		args := []string{expected.GetMeta().GetName()}
		args = append(args, tc.Args...)
		out, err := clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdUpdateDeployment(), args)
		if tc.Err == nil {
			require.NoError(t, err)

			var resp sdk.TxResponse
			require.NoError(t, nc.Ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, clitestutil.CheckTxCode(nc.Net, nc.Ctx, resp.TxHash, tc.Code))

			flagFilterBy := fmt.Sprintf(network.FlagPattern, cli.FlagFilterBy, "creator")
			flagFilterValue := fmt.Sprintf(network.FlagPattern, cli.FlagFilterValue, expected.GetMeta().Creator)

			out, err = clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdListDeployments(), []string{flagFilterBy, flagFilterValue})
			require.NoError(t, err)
			require.Contains(t, out.String(), expected.GetMeta().GetDomain())
			require.Contains(t, out.String(), expected.GetMeta().GetDescription())
			require.Contains(t, out.String(), expected.GetMeta().GetName())
			require.Contains(t, out.String(), expected.GetMeta().GetCreator())

		} else {
			require.Error(t, err)
			require.ErrorContains(t, err, tc.Err.Error())
		}
	})
}

func TestUpdateDeployment(t *testing.T) {
	nc := network.Setup(t)
	commonFlags := network.SetupTxCommonFlags(t, nc)

	testUpdateDomain(t, nc, commonFlags)
	testUpdateDescription(t, nc, commonFlags)
	testUpdateAll(t, nc, commonFlags)
}

func createDeployment(t *testing.T, nc *network.Context, commonFlags []string) *types.Deployment {
	archive, err := sample.CreateTempArchive("index.html", sample.HelloWorldHTMLBody)
	require.NoError(t, err)
	defer os.RemoveAll(archive.Name())

	deployment := &types.Deployment{
		Meta:    sample.CreateMetaWithAddr(nc.Val.Address.String(), 0),
		Dataset: sample.CreateDatasetFromStrings([]string{"index.html"}),
	}

	args := append([]string{deployment.GetMeta().GetName(), archive.Name()}, commonFlags...)
	_, err = clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdCreateDeployment(), args)
	require.NoError(t, err)
	return deployment
}

func testUpdateDomain(t *testing.T, nc *network.Context, commonFlags []string) {
	obj := createDeployment(t, nc, commonFlags)
	expected := obj
	expected.Meta.Domain = "foobar"
	flagDomain := fmt.Sprintf(network.FlagPattern, cli.FlagDomain, "foobar")

	runUpdateTxTest(t, nc, &network.TxTestCase{
		Name: "test update domain",
		Args: append([]string{flagDomain}, commonFlags...),
	}, expected)
}

func testUpdateDescription(t *testing.T, nc *network.Context, commonFlags []string) {
	obj := createDeployment(t, nc, commonFlags)
	expected := obj
	expected.Meta.Description = "hey ho"
	flagDescription := fmt.Sprintf(network.FlagPattern, cli.FlagDescription, "hey ho")

	runUpdateTxTest(t, nc, &network.TxTestCase{
		Name: "test update description",
		Args: append([]string{flagDescription}, commonFlags...),
	}, expected)
}

func testUpdateAll(t *testing.T, nc *network.Context, commonFlags []string) {
	obj := createDeployment(t, nc, commonFlags)
	expected := obj
	expected.Meta.Description = "new desc"
	expected.Meta.Domain = "barfoo"
	flagDescription := fmt.Sprintf(network.FlagPattern, cli.FlagDescription, "new desc")
	flagDomain := fmt.Sprintf(network.FlagPattern, cli.FlagDomain, "barfoo")

	runUpdateTxTest(t, nc, &network.TxTestCase{
		Name: "test update all",
		Args: append([]string{flagDescription, flagDomain}, commonFlags...),
	}, expected)
}

// TODO: Test update payload
