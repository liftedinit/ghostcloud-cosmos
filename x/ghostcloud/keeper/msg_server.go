package keeper

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"ghostcloud/x/ghostcloud/types"
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

func handlePayload(payload *types.Payload) (*types.Dataset, error) {
	if archive := payload.GetArchive(); archive != nil {
		return datasetFromArchive(archive)
	} else if dataset := payload.GetDataset(); dataset != nil {
		return dataset, nil
	}

	panic("invalid payload")
}

func validateMeta(meta *types.Meta) error {
	if meta.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if meta.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}
	return nil
}

func validatePayload(payload *types.Payload) error {
	if payload == nil {
		return fmt.Errorf("payload cannot be empty")
	}
	return nil
}
