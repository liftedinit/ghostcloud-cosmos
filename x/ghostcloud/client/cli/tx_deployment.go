package cli

import (
	"fmt"
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func ReadWebsiteRoot(path string) ([]*types.File, error) {
	// Walk through the directory and process each file
	var files []*types.File
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files = append(files, &types.File{
			Name:    info.Name(),
			Content: fileBytes,
		})

		if err != nil {
			fmt.Printf("Error encoding file %q: %v\n", path, err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func CmdCreateDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-deployment [name] [description] [domain] [memo] [website-root]",
		Short: "Create a new deployment",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexName := args[0]

			// Get value arguments
			argDescription := args[1]
			argDomain := args[2]
			argMemo := args[3]
			argWebsiteRoot := args[4]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			websiteFiles, err := ReadWebsiteRoot(argWebsiteRoot)
			if err != nil {
				return err
			}

			meta := types.Meta{
				Name:        indexName,
				Description: argDescription,
				Domain:      argDomain,
				Memo:        argMemo,
			}

			msg := types.NewMsgCreateDeployment(
				clientCtx.GetFromAddress().String(),
				&meta,
				websiteFiles,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// TODO: Create commands to update meta/files only

func CmdUpdateDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-deployment [name] [description] [domain] [memo] [website-root]",
		Short: "Update a deployment",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexName := args[0]

			// Get value arguments
			argDescription := args[1]
			argDomain := args[2]
			argMemo := args[3]
			argWebsiteRoot := args[4]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			websiteFiles, err := ReadWebsiteRoot(argWebsiteRoot)
			if err != nil {
				return err
			}

			meta := types.Meta{
				Name:        indexName,
				Description: argDescription,
				Domain:      argDomain,
				Memo:        argMemo,
			}

			msg := types.NewMsgUpdateDeployment(
				clientCtx.GetFromAddress().String(),
				&meta,
				websiteFiles,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-deployment [name]",
		Short: "Delete a deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexName := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteDeployment(
				clientCtx.GetFromAddress().String(),
				indexName,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
