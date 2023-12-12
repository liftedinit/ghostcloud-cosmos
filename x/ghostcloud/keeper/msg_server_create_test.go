package keeper_test

import (
	"fmt"
	"strings"
	"testing"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/keeper"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func testDeploymentMsgServerCreate(t *testing.T, k *keeper.Keeper, ctx sdk.Context, tc keepertest.MsgServerTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		srv := keeper.NewMsgServerImpl(*k)
		wctx := sdk.WrapSDKContext(ctx)
		require.Equal(t, len(tc.Metas), len(tc.Payloads))

		for i := 0; i < len(tc.Metas); i++ {
			expected := &types.MsgCreateDeploymentRequest{
				Meta:    tc.Metas[i],
				Payload: tc.Payloads[i],
			}
			_, err := srv.CreateDeployment(wctx, expected)
			if tc.Err == nil {
				require.NoError(t, err)

				creator, err := sdk.AccAddressFromBech32(expected.Meta.GetCreator())
				require.NoError(t, err)
				retrievedMeta, found := k.GetMeta(ctx, creator, expected.Meta.GetName())
				require.True(t, found)
				require.Equal(t, expected.Meta, &retrievedMeta)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.Err.Error())
			}
		}
	})
}

func testDeploymentMsgServerCreateValidDataset(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(keepertest.NUM_DEPLOYMENT, keepertest.DATASET_SIZE)
	require.Len(t, meta, keepertest.NUM_DEPLOYMENT)
	require.Len(t, payload, keepertest.NUM_DEPLOYMENT)
	tc := keepertest.MsgServerTestCase{
		Name:     "valid_d",
		Metas:    meta,
		Payloads: payload,
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgServerCreateDatasetNoIndex(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloads(keepertest.NUM_DEPLOYMENT, keepertest.DATASET_SIZE)
	require.Len(t, meta, keepertest.NUM_DEPLOYMENT)
	require.Len(t, payload, keepertest.NUM_DEPLOYMENT)
	tc := keepertest.MsgServerTestCase{
		Name:     "d_no_index",
		Metas:    meta,
		Payloads: payload,
		Err:      fmt.Errorf(types.IndexHtmlNotFound),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgServerCreateValidArchive(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNArchivePayloads(keepertest.NUM_DEPLOYMENT)
	require.Len(t, meta, keepertest.NUM_DEPLOYMENT)
	require.Len(t, payload, keepertest.NUM_DEPLOYMENT)
	tc := keepertest.MsgServerTestCase{
		Name:     "valid_a",
		Metas:    meta,
		Payloads: payload,
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgServerCreateArchiveTooBig(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	params := k.GetParams(ctx)
	meta, payload := sample.CreateRandomArchivePayload(1, params.MaxPayloadSize+1, "index.html")
	tc := keepertest.MsgServerTestCase{
		Name:     "a_too_big",
		Metas:    []*types.Meta{meta},
		Payloads: []*types.Payload{payload},
		Err:      fmt.Errorf("payload is too big"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgServerArchiveNoIndex(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateRandomArchivePayload(1, 1, "foobar.html")
	tc := keepertest.MsgServerTestCase{
		Name:     "a_no_index",
		Metas:    []*types.Meta{meta},
		Payloads: []*types.Payload{payload},
		Err:      fmt.Errorf(types.IndexHtmlNotFound),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgServerNameTooLong(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	params := k.GetParams(ctx)
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	meta[0].Name = strings.Repeat("a", int(params.MaxNameSize+1))
	tc := keepertest.MsgServerTestCase{
		Name:     "name_too_long",
		Metas:    meta,
		Payloads: payload,
		Err:      fmt.Errorf(types.NameTooLong, meta[0].Name),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerEmptyName(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	meta[0].Name = ""
	tc := keepertest.MsgServerTestCase{
		Name:     "empty_name",
		Metas:    meta,
		Payloads: payload,
		Err:      fmt.Errorf(types.NameShouldNotBeEmpty),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerEmptyCreator(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	meta[0].Creator = ""
	tc := keepertest.MsgServerTestCase{
		Name:     "empty_creator",
		Metas:    meta,
		Payloads: payload,
		Err:      fmt.Errorf(types.CreatorShouldNotBeEmpty),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerDescriptionTooLong(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	params := k.GetParams(ctx)
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	meta[0].Description = strings.Repeat("a", int(params.MaxDescriptionSize+1))
	tc := keepertest.MsgServerTestCase{
		Name:     "description_too_long",
		Metas:    meta,
		Payloads: payload,
		Err:      fmt.Errorf(types.DescriptionTooLong, meta[0].Description),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerEmptyPayload(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	payload[0] = nil
	tc := keepertest.MsgServerTestCase{
		Name:     "empty_payload",
		Metas:    meta,
		Payloads: payload,
		Err:      fmt.Errorf(types.PayloadIsRequired),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerEmptyArchivePayload(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNArchivePayloads(1)
	emptyArchivePayload := &types.Payload{
		PayloadOption: &types.Payload_Archive{
			Archive: nil,
		},
	}
	payload[0] = emptyArchivePayload
	tc := keepertest.MsgServerTestCase{
		Name:     "empty_archive_payload",
		Metas:    meta,
		Payloads: payload,
		Err:      fmt.Errorf("archive cannot be nil"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerEmptyDatasetPayload(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	emptyDatasetPayload := &types.Payload{
		PayloadOption: &types.Payload_Dataset{
			Dataset: nil,
		},
	}
	payload[0] = emptyDatasetPayload
	tc := keepertest.MsgServerTestCase{
		Name:     "empty_dataset_payload",
		Metas:    meta,
		Payloads: payload,
		Err:      fmt.Errorf("dataset cannot be nil"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerArchiveBombPayload(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	params := k.GetParams(ctx)
	meta, payload := sample.CreateBombArchivePayload(1, params.MaxUncompressedSize+1, "index.html")
	tc := keepertest.MsgServerTestCase{
		Name:     "a_bomb",
		Metas:    []*types.Meta{meta},
		Payloads: []*types.Payload{payload},
		Err:      fmt.Errorf(types.UncompressedSizeTooBig, params.MaxUncompressedSize+1, params.MaxUncompressedSize),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerInvalidArchiveType(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta := sample.CreateMeta(0)
	payload := &types.Payload{
		PayloadOption: &types.Payload_Archive{
			Archive: &types.Archive{Type: 123, Content: sample.CreateZip("index.html", "foobar")},
		},
	}
	tc := keepertest.MsgServerTestCase{
		Name:     "a_invalid_type",
		Metas:    []*types.Meta{meta},
		Payloads: []*types.Payload{payload},
		Err:      fmt.Errorf("unsupported archive type"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerUnsupportedPayloadType(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta := sample.CreateMeta(0)
	payload := &types.Payload{}
	tc := keepertest.MsgServerTestCase{
		Name:     "unsupported_payload_type",
		Metas:    []*types.Meta{meta},
		Payloads: []*types.Payload{payload},
		Err:      fmt.Errorf("unsupported payload type"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerNoMeta(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	_, payload := sample.CreateDatasetPayload(1, 1)
	tc := keepertest.MsgServerTestCase{
		Name:     "no_meta",
		Metas:    []*types.Meta{nil},
		Payloads: []*types.Payload{payload},
		Err:      fmt.Errorf("meta is required"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerInvalidCreator(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateDatasetPayloadWithIndexHtml(1, 1)
	meta.Creator = "invalid creator address"
	tc := keepertest.MsgServerTestCase{
		Name:     "invalid_creator",
		Metas:    []*types.Meta{meta},
		Payloads: []*types.Payload{payload},
		Err:      fmt.Errorf("invalid creator address"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerIndexAlreadySet(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	addr := sdk.MustAccAddressFromBech32(meta[0].Creator)
	k.SetDeployment(ctx, addr, meta[0], payload[0].GetDataset())
	tc := keepertest.MsgServerTestCase{
		Name:     "index_already_set",
		Metas:    []*types.Meta{meta[0]},
		Payloads: []*types.Payload{payload[0]},
		Err:      fmt.Errorf("index already set"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerNameHasWhitespace(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	meta[0].Name = "name with whitespace"
	tc := keepertest.MsgServerTestCase{
		Name:     "name_has_whitespace",
		Metas:    []*types.Meta{meta[0]},
		Payloads: []*types.Payload{payload[0]},
		Err:      fmt.Errorf("name should not contain whitespace"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerNameAsciiOnly(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	meta[0].Name = "™"
	tc := keepertest.MsgServerTestCase{
		Name:     "name_ascii_only",
		Metas:    []*types.Meta{meta[0]},
		Payloads: []*types.Payload{payload[0]},
		Err:      fmt.Errorf("name should contain ascii characters only"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func testDeploymentMsgCreateServerInvalidDomain(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	meta, payload := sample.CreateNDatasetPayloadsWithIndexHtml(1, keepertest.DATASET_SIZE)
	meta[0].Domain = "invalid domain"
	tc := keepertest.MsgServerTestCase{
		Name:     "invalid_domain",
		Metas:    []*types.Meta{meta[0]},
		Payloads: []*types.Payload{payload[0]},
		Err:      fmt.Errorf("invalid domain"),
	}
	testDeploymentMsgServerCreate(t, k, ctx, tc)
}

func TestDeploymentMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)

	testDeploymentMsgServerCreateValidDataset(t, k, ctx)
	testDeploymentMsgServerCreateDatasetNoIndex(t, k, ctx)
	testDeploymentMsgServerCreateValidArchive(t, k, ctx)
	testDeploymentMsgServerCreateArchiveTooBig(t, k, ctx)
	testDeploymentMsgServerArchiveNoIndex(t, k, ctx)
	testDeploymentMsgServerNameTooLong(t, k, ctx)
	testDeploymentMsgCreateServerEmptyName(t, k, ctx)
	testDeploymentMsgCreateServerEmptyCreator(t, k, ctx)
	testDeploymentMsgCreateServerDescriptionTooLong(t, k, ctx)
	testDeploymentMsgCreateServerEmptyPayload(t, k, ctx)
	testDeploymentMsgCreateServerEmptyArchivePayload(t, k, ctx)
	testDeploymentMsgCreateServerEmptyDatasetPayload(t, k, ctx)
	testDeploymentMsgCreateServerArchiveBombPayload(t, k, ctx)
	testDeploymentMsgCreateServerInvalidArchiveType(t, k, ctx)
	testDeploymentMsgCreateServerUnsupportedPayloadType(t, k, ctx)
	testDeploymentMsgCreateServerNoMeta(t, k, ctx)
	testDeploymentMsgCreateServerInvalidCreator(t, k, ctx)
	testDeploymentMsgCreateServerIndexAlreadySet(t, k, ctx)
	testDeploymentMsgCreateServerNameHasWhitespace(t, k, ctx)
	testDeploymentMsgCreateServerNameAsciiOnly(t, k, ctx)
	testDeploymentMsgCreateServerInvalidDomain(t, k, ctx)
}
