package cli_test

import (
	"context"
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

const (
	NewDescription = "new description"
	NewDomain      = "newdomain"
	NewContent     = "<h1>new content</h1>"
	IndexHTML      = "index.html"
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

			queryClient := types.NewQueryClient(nc.Val.ClientCtx)
			response, qerr := queryClient.Content(context.Background(), &types.QueryContentRequest{
				Creator: expected.GetMeta().GetCreator(),
				Name:    expected.GetMeta().GetName(),
				Path:    IndexHTML,
			})
			require.NoError(t, qerr)
			require.Equal(t, response.Content, expected.GetDataset().GetItems()[0].Content.Content)
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
	testUpdatePayload(t, nc, commonFlags)
	testUpdateAll(t, nc, commonFlags)
}

func setNewDescription(expected *types.Deployment) (*types.Deployment, string) {
	expected.Meta.Description = NewDescription
	return expected, fmt.Sprintf(network.FlagPattern, cli.FlagDescription, NewDescription)
}

func setNewDomain(expected *types.Deployment) (*types.Deployment, string) {
	expected.Meta.Domain = NewDomain
	return expected, fmt.Sprintf(network.FlagPattern, cli.FlagDomain, NewDomain)
}

func setNewPayload(expected *types.Deployment, newArchivePath string) (*types.Deployment, string) {
	expected.Dataset.Items[0].Content = &types.ItemContent{Content: []byte(NewContent)}
	return expected, fmt.Sprintf(network.FlagPattern, cli.FlagWebsitePayload, newArchivePath)
}

func createDeployment(t *testing.T, nc *network.Context, id int, commonFlags []string) *types.Deployment {
	archive, err := sample.CreateTempArchive(IndexHTML, sample.HelloWorldHTMLBody)
	require.NoError(t, err)
	defer os.RemoveAll(archive.Name())

	deployment := &types.Deployment{
		Meta: sample.CreateMetaWithAddr(nc.Val.Address.String(), id),
		Dataset: &types.Dataset{
			Items: []*types.Item{
				{
					Meta: &types.ItemMeta{Path: IndexHTML},
					Content: &types.ItemContent{
						Content: []byte(sample.HelloWorldHTMLBody),
					},
				},
			},
		},
	}

	args := append([]string{deployment.GetMeta().GetName(), archive.Name()}, commonFlags...)
	_, err = clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdCreateDeployment(), args)
	require.NoError(t, err)
	return deployment
}

func testUpdateDomain(t *testing.T, nc *network.Context, commonFlags []string) {
	expected := createDeployment(t, nc, 0, commonFlags)
	expected, flagDomain := setNewDomain(expected)

	runUpdateTxTest(t, nc, &network.TxTestCase{
		Name: "test update domain",
		Args: append([]string{flagDomain}, commonFlags...),
	}, expected)
}

func testUpdateDescription(t *testing.T, nc *network.Context, commonFlags []string) {
	expected := createDeployment(t, nc, 1, commonFlags)
	expected, flagDescription := setNewDescription(expected)

	runUpdateTxTest(t, nc, &network.TxTestCase{
		Name: "test update description",
		Args: append([]string{flagDescription}, commonFlags...),
	}, expected)
}

func testUpdatePayload(t *testing.T, nc *network.Context, commonFlags []string) {
	newArchive, err := sample.CreateTempArchive(IndexHTML, NewContent)
	require.NoError(t, err)
	defer os.RemoveAll(newArchive.Name())

	expected := createDeployment(t, nc, 2, commonFlags)
	expected, flagWebsitePayload := setNewPayload(expected, newArchive.Name())

	runUpdateTxTest(t, nc, &network.TxTestCase{
		Name: "test update payload",
		Args: append([]string{flagWebsitePayload}, commonFlags...),
	}, expected)
}

func testUpdateAll(t *testing.T, nc *network.Context, commonFlags []string) {
	newArchive, err := sample.CreateTempArchive(IndexHTML, NewContent)
	require.NoError(t, err)
	defer os.RemoveAll(newArchive.Name())

	expected := createDeployment(t, nc, 3, commonFlags)
	expected, flagDescription := setNewDescription(expected)
	expected, flagDomain := setNewDomain(expected)
	expected, flagWebsitePayload := setNewPayload(expected, newArchive.Name())

	runUpdateTxTest(t, nc, &network.TxTestCase{
		Name: "test update all",
		Args: append([]string{flagDescription, flagDomain, flagWebsitePayload}, commonFlags...),
	}, expected)
}
