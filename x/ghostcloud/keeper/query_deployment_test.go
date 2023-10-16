package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/testutil/nullify"
	"ghostcloud/x/ghostcloud/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestDeploymentQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.GhostcloudKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDeployment(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetDeploymentRequest
		response *types.QueryGetDeploymentResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDeploymentRequest{
				Name: msgs[0].Name,
			},
			response: &types.QueryGetDeploymentResponse{Deployment: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetDeploymentRequest{
				Name: msgs[1].Name,
			},
			response: &types.QueryGetDeploymentResponse{Deployment: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDeploymentRequest{
				Name: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Deployment(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestDeploymentQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.GhostcloudKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDeployment(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDeploymentRequest {
		return &types.QueryAllDeploymentRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DeploymentAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deployment), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Deployment),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DeploymentAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deployment), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Deployment),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.DeploymentAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Deployment),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.DeploymentAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
