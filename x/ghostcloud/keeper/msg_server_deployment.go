package keeper

import (
	"archive/zip"
	"bytes"
	"context"
	"io"

	"ghostcloud/x/ghostcloud/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const IdxNotSet = "index not set"
const InvalidCreatorAddr = "invalid creator address (%s)"

func (k msgServer) CreateDeployment(goCtx context.Context, msg *types.MsgCreateDeployment) (*types.MsgCreateDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
	}

	// Check if the value already exists
	_, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Meta.Name,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var deployment = types.Deployment{
		Creator: msg.Creator,
		Meta:    msg.Meta,
		Files:   msg.Files,
	}

	k.SetDeployment(
		ctx,
		deployment,
	)
	return &types.MsgCreateDeploymentResponse{}, nil
}

// filesFromArchiveBytes returns a list of files from an archive byte slice.
func filesFromArchiveBytes(archive []byte) ([]*types.File, error) {
	r := bytes.NewReader(archive)
	zipReader, err := zip.NewReader(r, int64(len(archive)))
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

	return files, nil
}

func (k msgServer) CreateDeploymentArchive(goCtx context.Context, msg *types.MsgCreateDeploymentArchive) (*types.MsgCreateDeploymentArchiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
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

	files, err := filesFromArchiveBytes(msg.WebsiteArchive)
	if err != nil {
		return nil, err
	}

	var deployment = types.Deployment{
		Creator: msg.Creator,
		Meta:    msg.Meta,
		Files:   files,
	}

	k.SetDeployment(ctx, deployment)

	return &types.MsgCreateDeploymentArchiveResponse{}, nil
}

func (k msgServer) UpdateDeployment(goCtx context.Context, msg *types.MsgUpdateDeployment) (*types.MsgUpdateDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
	}

	// Check if the value exists
	valFound, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Meta.Name,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, IdxNotSet)
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var deployment = types.Deployment{
		Creator: msg.Creator,
		Meta:    msg.Meta,
		Files:   msg.Files,
	}

	k.SetDeployment(ctx, deployment)

	return &types.MsgUpdateDeploymentResponse{}, nil
}

func (k msgServer) UpdateDeploymentMeta(goCtx context.Context, msg *types.MsgUpdateDeploymentMeta) (*types.MsgUpdateDeploymentMetaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
	}

	deployment, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Meta.Name,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, IdxNotSet)
	}

	deployment.Meta = msg.Meta

	k.SetDeployment(ctx, deployment)

	return &types.MsgUpdateDeploymentMetaResponse{}, nil
}

func (k msgServer) DeleteDeployment(goCtx context.Context, msg *types.MsgDeleteDeployment) (*types.MsgDeleteDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
	}

	// Check if the value exists
	valFound, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Name,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, IdxNotSet)
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveDeployment(
		ctx,
		addr,
		msg.Name,
	)

	return &types.MsgDeleteDeploymentResponse{}, nil
}
