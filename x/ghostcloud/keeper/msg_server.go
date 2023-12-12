package keeper

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"ghostcloud/x/ghostcloud/types"

	"github.com/asaskevich/govalidator"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

const InvalidCreatorAddr = "invalid creator address: %s"

func datasetFromZip(content []byte) (*types.Dataset, error) {
	r := bytes.NewReader(content)
	zipReader, err := zip.NewReader(r, int64(len(content)))
	if err != nil {
		return nil, fmt.Errorf("zip reader error: %w", err)
	}

	items := make([]*types.Item, 0, len(zipReader.File))
	for _, file := range zipReader.File {
		ferr := func(f *zip.File) error {
			rc, oerr := file.Open()
			if oerr != nil {
				return fmt.Errorf("error opening file: %w", oerr)
			}
			defer rc.Close()

			content, rerr := io.ReadAll(rc)
			if rerr != nil {
				return fmt.Errorf("error reading file: %w", rerr)
			}

			items = append(items, &types.Item{
				Meta:    &types.ItemMeta{Path: file.Name},
				Content: &types.ItemContent{Content: content},
			})
			return nil
		}(file)
		if ferr != nil {
			return nil, fmt.Errorf("error processing file: %w", ferr)
		}
	}

	return &types.Dataset{
		Items: items,
	}, nil
}
func datasetFromArchive(archive *types.Archive) (*types.Dataset, error) {
	switch archive.Type {
	case types.ArchiveType_Zip:
		return datasetFromZip(archive.Content)
	default:
		return nil, fmt.Errorf("unsupported archive type: %s", archive.Type)
	}
}

func HandlePayload(payload *types.Payload) (*types.Dataset, error) {
	if archive := payload.GetArchive(); archive != nil {
		return datasetFromArchive(archive)
	} else if dataset := payload.GetDataset(); dataset != nil {
		return dataset, nil
	}

	return nil, fmt.Errorf("unsupported payload type")
}

func validateCreator(creator string) error {
	if creator == "" {
		return fmt.Errorf(types.CreatorShouldNotBeEmpty)
	}

	return nil
}

func validateName(name string, maxNameSize int64) error {
	if name == "" {
		return fmt.Errorf(types.NameShouldNotBeEmpty)
	}
	if govalidator.HasWhitespace(name) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, types.NameShouldNotContainWhitespace, name)
	}
	if !govalidator.IsASCII(name) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, types.NameShouldContainASCII, name)
	}
	if int64(len(name)) > maxNameSize {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, types.NameTooLong, name)
	}

	return nil
}

func validateDescription(description string, maxDescriptionSize int64) error {
	if int64(len(description)) > maxDescriptionSize {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, types.DescriptionTooLong, description)
	}

	return nil
}

func validateDomain(domain string) error {
	if domain != "" && !govalidator.IsDNSName(domain) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, types.InvalidDomain, domain)
	}
	return nil
}

func validateMeta(meta *types.Meta, params types.Params) error {
	if meta == nil {
		return fmt.Errorf(types.MetaIsRequired)
	}
	if err := validateName(meta.Name, params.MaxNameSize); err != nil {
		return err
	}
	if err := validateCreator(meta.Creator); err != nil {
		return err
	}
	if err := validateDescription(meta.Description, params.MaxDescriptionSize); err != nil {
		return err
	}
	if err := validateDomain(meta.Domain); err != nil {
		return err
	}

	return nil
}

func validatePayload(payload *types.Payload, params types.Params) error {
	if int64(payload.Size()) > params.MaxPayloadSize {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, types.PayloadTooBig, payload.Size(), params.MaxPayloadSize)
	}

	switch payload.GetPayloadOption().(type) {
	case *types.Payload_Archive:
		archive := payload.GetArchive()
		if archive == nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "archive cannot be nil")
		}
		if err := verifyArchiveContent(archive.Content, params.MaxUncompressedSize); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
		}
	case *types.Payload_Dataset:
		dataset := payload.GetDataset()
		if dataset == nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "dataset cannot be nil")
		}
		if err := verifyDatasetContent(dataset); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}
	return nil
}

func verifyArchiveContent(archive []byte, maxUncompressedSize uint64) error {
	r := bytes.NewReader(archive)
	zipReader, err := zip.NewReader(r, int64(len(archive)))
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var totalUncompressedSize uint64
	var indexFound bool
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		totalUncompressedSize += file.UncompressedSize64
		if totalUncompressedSize > maxUncompressedSize {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, types.UncompressedSizeTooBig, totalUncompressedSize, maxUncompressedSize)
		}

		if file.Name == "index.html" {
			indexFound = true
		}
	}

	if !indexFound {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, types.IndexHtmlNotFound)
	}

	return nil
}

func verifyDatasetContent(dataset *types.Dataset) error {
	var indexFound bool
	for _, item := range dataset.Items {
		if item.Meta.Path == "index.html" {
			indexFound = true
		}
	}

	if !indexFound {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, types.IndexHtmlNotFound)
	}

	return nil
}
