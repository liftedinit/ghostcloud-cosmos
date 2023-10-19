package keeper

import (
	"archive/zip"
	"bytes"
	"context"
	"ghostcloud/x/ghostcloud/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"io"
)

func (k msgServer) CreateDeploymentArchive(goCtx context.Context, msg *types.MsgCreateDeploymentArchive) (*types.MsgCreateDeploymentArchiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Check if the value exists
	_, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Meta.Name,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	// TODO: Refactor
	r := bytes.NewReader(msg.WebsiteArchive)
	zipReader, err := zip.NewReader(r, int64(len(msg.WebsiteArchive)))
	if err != nil {
		return nil, err
	}

	files := make([]*types.File, 0)
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(rc)
		if err != nil {
			cerr := rc.Close()
			if cerr != nil {
				return nil, cerr
			}
			return nil, err
		}
		err = rc.Close()
		if err != nil {
			return nil, err
		}

		files = append(files, &types.File{
			Meta: &types.FileMeta{
				Name: file.Name,
			},
			Content: content,
		})
	}

	var deployment = types.Deployment{
		Creator: msg.Creator,
		Meta:    msg.Meta,
		Files:   files,
	}

	k.SetDeployment(ctx, deployment)

	return &types.MsgCreateDeploymentArchiveResponse{}, nil
}
