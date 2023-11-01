package cli

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"ghostcloud/x/ghostcloud/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/spf13/cobra"
)

const (
	FlagDescription  = "description"
	FlagDomain       = "domain"
	zipArchiveSuffix = ".zip"
)

func addDeploymentFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.String(FlagDescription, "", "Description of the deployment")
	f.String(FlagDomain, "", "Custom domain of the deployment")
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	cmd.AddCommand(CmdCreateDeployment())

	return cmd
}

// isDir Check if a path is a directory. Panics if the path does not exist.
func isDir(path string) bool {
	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("File does not exist: %s", path)
	}
	return info.IsDir()
}

func loadArchive(path string) []byte {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatalf("unable to stat website archive: %v", err)
	}
	if fileInfo.Size() > types.DefaultMaxArchiveSize {
		log.Fatalf("Website archive is too big: %d > %d", fileInfo.Size(), types.DefaultMaxArchiveSize)
	}

	// Read website archive
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read website archive: %v", err)
	}

	r := bytes.NewReader(data)
	zipReader, err := zip.NewReader(r, int64(len(data)))
	if err != nil {
		log.Fatalf("unable to create website archive reader: %v", err)
	}

	found := false
	for _, f := range zipReader.File {
		if f.Name == "index.html" {
			found = true
		}
	}

	if !found {
		log.Fatal("Website archive does not contain `index.html` at its root")
	}

	return data
}

func loadFolder(path string) []*types.Item {
	// Walk through the directory and process each file
	var items []*types.Item
	werr := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		content, rerr := os.ReadFile(path)
		if rerr != nil {
			log.Fatalf("unable to read file: %v", rerr)
		}
		items = append(items, &types.Item{
			Meta: &types.ItemMeta{
				Path: info.Name(),
			},
			Content: &types.ItemContent{
				Content: content,
			},
		})

		return nil
	})
	if werr != nil {
		log.Fatalf("unable to walk through website folder: %v", werr)
	}

	return items
}

func createArchivePayload(path string) *types.Payload {
	data := loadArchive(path)
	return &types.Payload{
		PayloadOption: &types.Payload_Archive{
			Archive: &types.Archive{
				Type:    types.ArchiveType_Zip,
				Content: data,
			},
		},
	}
}

func createDatasetPayload(path string) *types.Payload {
	data := loadFolder(path)
	return &types.Payload{
		PayloadOption: &types.Payload_Dataset{
			Dataset: &types.Dataset{
				Items: data,
			},
		},
	}
}

func createPayload(path string) *types.Payload {
	if strings.HasSuffix(path, zipArchiveSuffix) {
		return createArchivePayload(path)
	} else if isDir(path) {
		return createDatasetPayload(path)
	}

	log.Fatal("Website payload must be a directory or a zip archive")
	panic("unreachable")
}

func CmdCreateDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create name website-payload",
		Short: "Create a new deployment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			argName := args[0]
			argWebsitePayload := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			payload := createPayload(argWebsitePayload)

			argDescription := cmd.Flag(FlagDescription).Value.String()
			argDomain := cmd.Flag(FlagDomain).Value.String()

			meta := types.Meta{
				Creator:     clientCtx.GetFromAddress().String(),
				Name:        argName,
				Description: argDescription,
				Domain:      argDomain,
			}

			msg := types.NewMsgCreateDeploymentRequest(
				&meta,
				payload,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addDeploymentFlags(cmd)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
