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

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	cosmosnetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
)

type NetworkContext struct {
	net *network.Network
	val *cosmosnetwork.Validator
	ctx client.Context
}

type TestCase struct {
	name  string
	valid bool
	args  []string
	err   error
	code  uint32
}

func setup(t *testing.T) (*NetworkContext, []string) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx
	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10))).String()),
	}
	return &NetworkContext{
		net: net,
		val: val,
		ctx: ctx,
	}, commonFlags
}

func run(t *testing.T, nc *NetworkContext, tc *TestCase) {
	t.Run(tc.name, func(t *testing.T) {
		require.NoError(t, nc.net.WaitForNextBlock())

		args := []string{tc.name}
		args = append(args, tc.args...)
		out, err := clitestutil.ExecTestCLICmd(nc.ctx, cli.CmdCreateDeployment(), args)
		if tc.valid {
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			var resp sdk.TxResponse
			require.NoError(t, nc.ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, clitestutil.CheckTxCode(nc.net, nc.ctx, resp.TxHash, tc.code))
		} else {
			require.Error(t, err)
			require.ErrorContains(t, err, tc.err.Error())
		}
	})
}

func TestCreateDeployment(t *testing.T) {
	nc, commonFlags := setup(t)

	testValidDataset(t, nc, commonFlags)
	testValidArchive(t, nc, commonFlags)
	testInvalidDatasetPath(t, nc, commonFlags)
	testInvalidArchivePath(t, nc, commonFlags)
	testArchiveTooBig(t, nc, commonFlags)
	testNoIndex(t, nc, commonFlags)
}

func testValidDataset(t *testing.T, nc *NetworkContext, commonFlags []string) {
	data, err := sample.CreateTempDataset()
	require.NoError(t, err)
	defer os.RemoveAll(data)

	run(t, nc, &TestCase{
		name:  "valid_dataset",
		args:  append([]string{data}, commonFlags...),
		valid: true,
	})
}

func testValidArchive(t *testing.T, nc *NetworkContext, commonFlags []string) {
	data, err := sample.CreateTempArchive("index.html")
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	run(t, nc, &TestCase{
		name:  "valid_archive",
		args:  append([]string{data.Name()}, commonFlags...),
		valid: true,
	})
}

func testInvalidDatasetPath(t *testing.T, nc *NetworkContext, commonFlags []string) {
	run(t, nc, &TestCase{
		name:  "invalid_dataset_path",
		args:  append([]string{"some-invalid-path"}, commonFlags...),
		err:   fmt.Errorf("unable to stat path"),
		valid: false,
	})
}

func testInvalidArchivePath(t *testing.T, nc *NetworkContext, commonFlags []string) {
	run(t, nc, &TestCase{
		name:  "invalid_archive_path",
		args:  append([]string{"some-invalid-path.zip"}, commonFlags...),
		err:   fmt.Errorf("unable to stat archive"),
		valid: false,
	})
}

func testArchiveTooBig(t *testing.T, nc *NetworkContext, commonFlags []string) {
	data, err := sample.CreateCustomFakeArchive(types.DefaultMaxArchiveSize + 1)
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	run(t, nc, &TestCase{
		name:  "archive_too_big",
		args:  append([]string{data.Name()}, commonFlags...),
		err:   fmt.Errorf("website archive is too big"),
		valid: false,
	})
}

func testNoIndex(t *testing.T, nc *NetworkContext, commonFlags []string) {
	data, err := sample.CreateTempArchive("foobar.html")
	require.NoError(t, err)
	defer data.Close()
	defer os.Remove(data.Name())

	run(t, nc, &TestCase{
		name:  "no_index",
		args:  append([]string{data.Name()}, commonFlags...),
		err:   fmt.Errorf("website archive does not contain `index.html` at its root"),
		valid: false,
	})
}
