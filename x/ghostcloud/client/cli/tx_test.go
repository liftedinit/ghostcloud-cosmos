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

type TxTestCase struct {
	name string
	args []string
	err  error
	code uint32
}

func runTxTest(t *testing.T, nc *network.Context, tc *TxTestCase) {
	t.Run(tc.name, func(t *testing.T) {
		require.NoError(t, nc.Net.WaitForNextBlock())

		args := []string{tc.name}
		args = append(args, tc.args...)
		out, err := clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdCreateDeployment(), args)
		if tc.err == nil {
			require.NoError(t, err)

			var resp sdk.TxResponse
			require.NoError(t, nc.Ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, clitestutil.CheckTxCode(nc.Net, nc.Ctx, resp.TxHash, tc.code))
		} else {
			require.Error(t, err)
			require.ErrorContains(t, err, tc.err.Error())
		}
	})
}

func TestCreateDeployment(t *testing.T) {
	nc := network.Setup(t)
	commonFlags := network.SetupTxCommonFlags(t, nc)

	testValidDataset(t, nc, commonFlags)
	testValidArchive(t, nc, commonFlags)
	testInvalidDatasetPath(t, nc, commonFlags)
	testInvalidArchivePath(t, nc, commonFlags)
	testArchiveTooBig(t, nc, commonFlags)
	testNoIndex(t, nc, commonFlags)
}

func testValidDataset(t *testing.T, nc *network.Context, commonFlags []string) {
	data, err := sample.CreateTempDataset()
	require.NoError(t, err)
	defer os.RemoveAll(data)

	runTxTest(t, nc, &TxTestCase{
		name: "valid_dataset",
		args: append([]string{data}, commonFlags...),
	})
}

func testValidArchive(t *testing.T, nc *network.Context, commonFlags []string) {
	data, err := sample.CreateTempArchive("index.html")
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	runTxTest(t, nc, &TxTestCase{
		name: "valid_archive",
		args: append([]string{data.Name()}, commonFlags...),
	})
}

func testInvalidDatasetPath(t *testing.T, nc *network.Context, commonFlags []string) {
	runTxTest(t, nc, &TxTestCase{
		name: "invalid_dataset_path",
		args: append([]string{"some-invalid-path"}, commonFlags...),
		err:  fmt.Errorf("unable to stat path"),
	})
}

func testInvalidArchivePath(t *testing.T, nc *network.Context, commonFlags []string) {
	runTxTest(t, nc, &TxTestCase{
		name: "invalid_archive_path",
		args: append([]string{"some-invalid-path.zip"}, commonFlags...),
		err:  fmt.Errorf("unable to stat archive"),
	})
}

func testArchiveTooBig(t *testing.T, nc *network.Context, commonFlags []string) {
	data, err := sample.CreateCustomFakeArchive(types.DefaultMaxArchiveSize + 1)
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	runTxTest(t, nc, &TxTestCase{
		name: "archive_too_big",
		args: append([]string{data.Name()}, commonFlags...),
		err:  fmt.Errorf("website archive is too big"),
	})
}

func testNoIndex(t *testing.T, nc *network.Context, commonFlags []string) {
	data, err := sample.CreateTempArchive("foobar.html")
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	runTxTest(t, nc, &TxTestCase{
		name: "no_index",
		args: append([]string{data.Name()}, commonFlags...),
		err:  fmt.Errorf("website archive does not contain `index.html` at its root"),
	})
}
