package cli_test

import (
	"fmt"
	"os"
	"testing"

	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/client/cli"
	"ghostcloud/x/ghostcloud/types"

	"ghostcloud/testutil/network"

	"github.com/stretchr/testify/require"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func runCreateTxTest(t *testing.T, nc *network.Context, tc *network.TxTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		require.NoError(t, nc.Net.WaitForNextBlock())

		args := []string{tc.Name}
		args = append(args, tc.Args...)
		out, err := clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdCreateDeployment(), args)
		if tc.Err == nil {
			require.NoError(t, err)

			var resp sdk.TxResponse
			require.NoError(t, nc.Ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, clitestutil.CheckTxCode(nc.Net, nc.Ctx, resp.TxHash, tc.Code))
		} else {
			require.Error(t, err)
			require.ErrorContains(t, err, tc.Err.Error())
		}
	})
}

func TestCreateDeployment(t *testing.T) {
	nc := network.Setup(t)
	commonFlags := network.SetupTxCommonFlags(t, nc)

	testCreateValidDataset(t, nc, commonFlags)
	testCreateValidArchive(t, nc, commonFlags)
	testCreateInvalidDatasetPath(t, nc, commonFlags)
	testCreateInvalidArchivePath(t, nc, commonFlags)
	testCreateArchiveTooBig(t, nc, commonFlags)
	testCreateNoIndex(t, nc, commonFlags)
}

func testCreateValidDataset(t *testing.T, nc *network.Context, commonFlags []string) {
	data, err := sample.CreateTempDataset()
	require.NoError(t, err)
	defer os.RemoveAll(data)

	runCreateTxTest(t, nc, &network.TxTestCase{
		Name: "valid_d",
		Args: append([]string{data}, commonFlags...),
	})
}

func testCreateValidArchive(t *testing.T, nc *network.Context, commonFlags []string) {
	data, err := sample.CreateTempArchive("index.html", sample.HelloWorldHTMLBody)
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	runCreateTxTest(t, nc, &network.TxTestCase{
		Name: "valid_a",
		Args: append([]string{data.Name()}, commonFlags...),
	})
}

func testCreateInvalidDatasetPath(t *testing.T, nc *network.Context, commonFlags []string) {
	runCreateTxTest(t, nc, &network.TxTestCase{
		Name: "invalid_d",
		Args: append([]string{"some-invalid-path"}, commonFlags...),
		Err:  fmt.Errorf("website payload does not exist"),
	})
}

func testCreateInvalidArchivePath(t *testing.T, nc *network.Context, commonFlags []string) {
	runCreateTxTest(t, nc, &network.TxTestCase{
		Name: "invalid_ah",
		Args: append([]string{"some-invalid-path.zip"}, commonFlags...),
		Err:  fmt.Errorf("website payload does not exist"),
	})
}

func testCreateArchiveTooBig(t *testing.T, nc *network.Context, commonFlags []string) {
	data, err := sample.CreateCustomFakeArchive(types.DefaultMaxArchiveSize + 1)
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	runCreateTxTest(t, nc, &network.TxTestCase{
		Name: "a_too_big",
		Args: append([]string{data.Name()}, commonFlags...),
		Err:  fmt.Errorf("website archive is too big"),
	})
}

func testCreateNoIndex(t *testing.T, nc *network.Context, commonFlags []string) {
	data, err := sample.CreateTempArchive("foobar.html", sample.HelloWorldHTMLBody)
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	runCreateTxTest(t, nc, &network.TxTestCase{
		Name: "no_index",
		Args: append([]string{data.Name()}, commonFlags...),
		Err:  fmt.Errorf("website archive does not contain `index.html` at its root"),
	})
}
