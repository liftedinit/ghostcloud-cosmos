package cli_test

import (
	"fmt"
	"ghostcloud/testutil/keeper"
	"ghostcloud/testutil/network"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/client/cli"
	"ghostcloud/x/ghostcloud/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
	"testing"
)

func networkWithDeploymentObjects(t *testing.T, n int) (*network.Network, []*types.Deployment) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{
		Params:      types.DefaultParams(),
		Deployments: sample.CreateNDeployments(n, keeper.DATASET_SIZE),
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.Deployments
}

func metasFromDeployments(deployments []*types.Deployment) []*types.Meta {
	metas := make([]*types.Meta, len(deployments))
	for i, deployment := range deployments {
		metas[i] = deployment.Meta
	}
	return metas
}

func TestShowDeployment(t *testing.T) {
	net, objs := networkWithDeploymentObjects(t, keeper.NUM_DEPLOYMENT)
	metas := metasFromDeployments(objs)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc  string
		args  []string
		err   error
		metas []*types.Meta
	}{
		{
			desc:  "found",
			args:  common,
			metas: metas,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDeployments(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryMetasResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Meta)
				require.ElementsMatch(t, tc.metas, resp.Meta)
			}
		})
	}
}
