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

func testDeploymentMsgServerUpdate(t *testing.T, k *keeper.Keeper, ctx sdk.Context, tc keepertest.MsgServerTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		srv := keeper.NewMsgServerImpl(*k)
		wctx := sdk.WrapSDKContext(ctx)

		for i := 0; i < len(tc.Metas); i++ {
			meta := tc.Metas[i]
			payload := tc.Payloads[i]
			_, err := srv.UpdateDeployment(wctx, &types.MsgUpdateDeploymentRequest{Meta: meta, Payload: payload})
			if tc.Err == nil {
				require.NoError(t, err)
				creator, err := sdk.AccAddressFromBech32(meta.GetCreator())
				require.NoError(t, err)
				storeMeta, found := k.GetMeta(ctx, creator, meta.GetName())
				require.True(t, found)
				require.Equal(t, meta, &storeMeta)

				storeDataset := k.GetDataset(ctx, creator, meta.GetName())
				require.True(t, found)
				switch payload.GetPayloadOption().(type) {
				case *types.Payload_Archive:
					dataset, err := keeper.HandlePayload(payload)
					require.NoError(t, err)
					require.Equal(t, dataset, storeDataset)
				case *types.Payload_Dataset:
					require.Equal(t, payload.GetDataset(), storeDataset)
				}
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.Err.Error())
			}
		}
	})
}

func testDeploymentMsgServerUpdateValidDataset(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	newMeta, newPayload := sample.CreateDatasetPayloadWithAddrAndIndexHtml(metas[0].Creator, 0, 1)
	tc := keepertest.MsgServerTestCase{
		Name:     "update_valid_dataset",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{newPayload},
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateValidArchive(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	newMeta, newPayload := sample.CreateArchivePayloadWithAddrAndIndexHtml(metas[0].Creator, 0)
	tc := keepertest.MsgServerTestCase{
		Name:     "update_valid_archive",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{newPayload},
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateNoMeta(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	tc := keepertest.MsgServerTestCase{
		Name:     "no_meta",
		Metas:    []*types.Meta{nil},
		Payloads: []*types.Payload{nil},
		Err:      fmt.Errorf("meta is required"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateEmptyName(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	newMeta := sample.CreateMeta(0)
	newMeta.Name = ""
	tc := keepertest.MsgServerTestCase{
		Name:     "update_empty_name",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{nil},
		Err:      fmt.Errorf(types.NameShouldNotBeEmpty),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateNameTooLong(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	newMeta := sample.CreateMeta(0)
	newMeta.Name = strings.Repeat("a", int(k.GetParams(ctx).MaxNameSize+1))
	tc := keepertest.MsgServerTestCase{
		Name:     "update_name_too_long",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{nil},
		Err:      fmt.Errorf(types.NameTooLong, newMeta.Name),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateDescriptionTooLong(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	newMeta := sample.CreateMeta(0)
	newMeta.Description = strings.Repeat("a", int(k.GetParams(ctx).MaxDescriptionSize+1))
	tc := keepertest.MsgServerTestCase{
		Name:     "update_description_too_long",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{nil},
		Err:      fmt.Errorf(types.DescriptionTooLong, newMeta.Description),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdatePayloadTooBig(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	params := k.GetParams(ctx)
	newMeta, newPayload := sample.CreateRandomArchivePayload(1, params.MaxPayloadSize+1, "index.html")
	tc := keepertest.MsgServerTestCase{
		Name:     "update_payload_too_big",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{newPayload},
		Err:      fmt.Errorf(types.PayloadTooBig, newPayload.Size(), k.GetParams(ctx).MaxPayloadSize),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateArchiveBomb(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	params := k.GetParams(ctx)
	newMeta, newPayload := sample.CreateBombArchivePayload(1, params.MaxUncompressedSize+1, "index.html")
	tc := keepertest.MsgServerTestCase{
		Name:     "update_archive_bomb",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{newPayload},
		Err:      fmt.Errorf(types.UncompressedSizeTooBig, params.MaxUncompressedSize+1, params.MaxUncompressedSize),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateEmptyArchive(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 3)
	newMeta := sample.CreateMeta(0)
	newPayload := &types.Payload{
		PayloadOption: &types.Payload_Archive{Archive: nil},
	}
	tc := keepertest.MsgServerTestCase{
		Name:     "update_empty_archive",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{newPayload},
		Err:      fmt.Errorf("archive cannot be nil"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateArchiveNoIndexHtml(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta, newPayload := sample.CreateRandomArchivePayload(0, 1, "foobar.html")
	tc := keepertest.MsgServerTestCase{
		Name:     "update_archive_no_index_html",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{newPayload},
		Err:      fmt.Errorf(types.IndexHtmlNotFound),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateEmptyDataset(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta := sample.CreateMeta(0)
	newPayload := &types.Payload{
		PayloadOption: &types.Payload_Dataset{Dataset: nil},
	}
	tc := keepertest.MsgServerTestCase{
		Name:     "update_empty_dataset",
		Metas:    []*types.Meta{newMeta},
		Payloads: []*types.Payload{newPayload},
		Err:      fmt.Errorf("dataset cannot be nil"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateDatasetNoIndexHtml(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta, newPayload := sample.CreateNDatasetPayloads(1, 1)
	tc := keepertest.MsgServerTestCase{
		Name:     "update_dataset_no_index_html",
		Metas:    newMeta,
		Payloads: newPayload,
		Err:      fmt.Errorf(types.IndexHtmlNotFound),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateEmptyCreator(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta, newPayload := sample.CreateNDatasetPayloads(1, 1)
	newMeta[0].Creator = ""
	tc := keepertest.MsgServerTestCase{
		Name:     "update_empty_creator",
		Metas:    newMeta,
		Payloads: newPayload,
		Err:      fmt.Errorf("creator should not be empty"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateInvalidCreator(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta, newPayload := sample.CreateNDatasetPayloadsWithIndexHtml(1, 1)
	newMeta[0].Creator = "invalid_creator"
	tc := keepertest.MsgServerTestCase{
		Name:     "update_invalid_creator",
		Metas:    newMeta,
		Payloads: newPayload,
		Err:      fmt.Errorf("invalid creator address"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateNonExisting(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta, newPayload := sample.CreateNDatasetPayloadsWithIndexHtml(1, 1)
	tc := keepertest.MsgServerTestCase{
		Name:     "update_non_existing",
		Metas:    newMeta,
		Payloads: newPayload,
		Err:      fmt.Errorf("unable to update a non-existing deployment"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateUnsupportedPayloadType(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newPayload := &types.Payload{}
	tc := keepertest.MsgServerTestCase{
		Name:     "update_unsupported_payload_type",
		Metas:    metas,
		Payloads: []*types.Payload{newPayload},
		Err:      fmt.Errorf("unsupported payload type"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateUnsupportedArchiveType(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newPayload := &types.Payload{
		PayloadOption: &types.Payload_Archive{
			Archive: &types.Archive{Type: 123, Content: sample.CreateZip("index.html", "foobar")},
		},
	}
	tc := keepertest.MsgServerTestCase{
		Name:     "update_unsupported_archive_type",
		Metas:    metas,
		Payloads: []*types.Payload{newPayload},
		Err:      fmt.Errorf("unsupported archive type"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateNameWithWhitespace(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta, _ := sample.CreateNDatasetPayloadsWithIndexHtml(1, 1)
	newMeta[0].Name = "invalid name"
	tc := keepertest.MsgServerTestCase{
		Name:     "update_name_with_whitespace",
		Metas:    newMeta,
		Payloads: []*types.Payload{nil},
		Err:      fmt.Errorf("name should not contain whitespace"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateNameAsciiOnly(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta, _ := sample.CreateNDatasetPayloadsWithIndexHtml(1, 1)
	newMeta[0].Name = "™"
	tc := keepertest.MsgServerTestCase{
		Name:     "update_name_ascii_only",
		Metas:    newMeta,
		Payloads: []*types.Payload{nil},
		Err:      fmt.Errorf("name should contain ascii characters only"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateInvalidDomain(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	newMeta, _ := sample.CreateNDatasetPayloadsWithIndexHtml(1, 1)
	newMeta[0].Domain = "invalid domain"
	tc := keepertest.MsgServerTestCase{
		Name:     "update_invalid_domain",
		Metas:    newMeta,
		Payloads: []*types.Payload{nil},
		Err:      fmt.Errorf("invalid domain"),
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateRemoveDomain(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	metas[0].Domain = ""
	tc := keepertest.MsgServerTestCase{
		Name:     "update_remove_domain",
		Metas:    metas,
		Payloads: []*types.Payload{nil},
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func testDeploymentMsgServerUpdateRemoveDescription(t *testing.T, k *keeper.Keeper, ctx sdk.Context) {
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 1, 1)
	metas[0].Description = ""
	tc := keepertest.MsgServerTestCase{
		Name:     "update_remove_description",
		Metas:    metas,
		Payloads: []*types.Payload{nil},
	}
	testDeploymentMsgServerUpdate(t, k, ctx, tc)
}

func TestDeploymentMsgServerUpdate(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)

	testDeploymentMsgServerUpdateValidDataset(t, k, ctx)
	testDeploymentMsgServerUpdateValidArchive(t, k, ctx)
	testDeploymentMsgServerUpdateNoMeta(t, k, ctx)
	testDeploymentMsgServerUpdateEmptyName(t, k, ctx)
	testDeploymentMsgServerUpdateNameTooLong(t, k, ctx)
	testDeploymentMsgServerUpdateDescriptionTooLong(t, k, ctx)
	testDeploymentMsgServerUpdatePayloadTooBig(t, k, ctx)
	testDeploymentMsgServerUpdateArchiveBomb(t, k, ctx)
	testDeploymentMsgServerUpdateArchiveNoIndexHtml(t, k, ctx)
	testDeploymentMsgServerUpdateEmptyArchive(t, k, ctx)
	testDeploymentMsgServerUpdateEmptyDataset(t, k, ctx)
	testDeploymentMsgServerUpdateDatasetNoIndexHtml(t, k, ctx)
	testDeploymentMsgServerUpdateEmptyCreator(t, k, ctx)
	testDeploymentMsgServerUpdateInvalidCreator(t, k, ctx)
	testDeploymentMsgServerUpdateNonExisting(t, k, ctx)
	testDeploymentMsgServerUpdateUnsupportedPayloadType(t, k, ctx)
	testDeploymentMsgServerUpdateUnsupportedArchiveType(t, k, ctx)
	testDeploymentMsgServerUpdateNameWithWhitespace(t, k, ctx)
	testDeploymentMsgServerUpdateNameAsciiOnly(t, k, ctx)
	testDeploymentMsgServerUpdateInvalidDomain(t, k, ctx)
	testDeploymentMsgServerUpdateRemoveDomain(t, k, ctx)
	testDeploymentMsgServerUpdateRemoveDescription(t, k, ctx)
}
