package cli_test

import (
	"testing"

	"ghostcloud/testutil/keeper"
	"ghostcloud/testutil/network"
	"ghostcloud/x/ghostcloud/client/cli"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

type QueryTestCase struct {
	name  string
	args  []string
	err   error
	metas []*types.Meta
}

func metasFromDeployments(deployments []*types.Deployment) []*types.Meta {
	metas := make([]*types.Meta, len(deployments))
	for i, deployment := range deployments {
		metas[i] = deployment.Meta
	}
	return metas
}

func runQueryTest(t *testing.T, nc *network.Context, tc *QueryTestCase) {
	t.Run(tc.name, func(t *testing.T) {
		var args []string
		args = append(args, tc.args...)
		out, err := clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdListDeployments(), args)
		if tc.err != nil {
			stat, ok := status.FromError(tc.err)
			require.True(t, ok)
			require.ErrorIs(t, stat.Err(), tc.err)
		} else {
			require.NoError(t, err)
			var resp types.QueryMetasResponse
			require.NoError(t, nc.Net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NotNil(t, resp.Meta)
			require.ElementsMatch(t, tc.metas, resp.Meta)
		}
	})
}

func testListDeployments(t *testing.T, nc *network.Context, commonFlags []string, objs []*types.Deployment) {
	metas := metasFromDeployments(objs)
	runQueryTest(t, nc, &QueryTestCase{
		name:  "found",
		args:  commonFlags,
		metas: metas,
	})
}

func TestQueries(t *testing.T) {
	nc, objs := network.SetupWithDeployments(t, keeper.NUM_DEPLOYMENT)
	commonFlags := network.SetupQueryCommonFlags(t)

	testListDeployments(t, nc, commonFlags, objs)
}
