package cli_test

import (
	"fmt"
	"ghostcloud/testutil/sample"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ghostcloud/testutil/network"
	"ghostcloud/testutil/nullify"
	"ghostcloud/x/ghostcloud/client/cli"
	"ghostcloud/x/ghostcloud/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithDeploymentObjects(t *testing.T, n int) (*network.Network, []types.Deployment) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		deployment := types.Deployment{
			Creator: sample.AccAddress(),
			Meta:    sample.GetDeploymentMeta(i),
			Files:   sample.GetDeploymentFiles(i),
		}
		nullify.Fill(&deployment)
		state.DeploymentList = append(state.DeploymentList, deployment)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.DeploymentList
}

func TestShowDeployment(t *testing.T) {
	net, objs := networkWithDeploymentObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc      string
		idName    string
		idCreator string

		args []string
		err  error
		obj  types.Deployment
	}{
		{
			desc:      "found",
			idCreator: objs[0].Creator,
			idName:    objs[0].Meta.Name,

			args: common,
			obj:  objs[0],
		},
		{
			desc:      "not found",
			idName:    strconv.Itoa(100000),
			idCreator: "B",

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idName,
				tc.idCreator,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowDeployment(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetDeploymentResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Deployment)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Deployment),
				)
			}
		})
	}
}

func TestListDeployment(t *testing.T) {
	net, objs := networkWithDeploymentObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDeployment(), args)
			require.NoError(t, err)
			var resp types.QueryAllDeploymentResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Deployment), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Deployment),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDeployment(), args)
			require.NoError(t, err)
			var resp types.QueryAllDeploymentResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Deployment), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Deployment),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDeployment(), args)
		require.NoError(t, err)
		var resp types.QueryAllDeploymentResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.Deployment),
		)
	})
}
