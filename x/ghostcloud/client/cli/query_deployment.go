package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"ghostcloud/x/ghostcloud/types"
)

// TODO: Add query command to retrieve the file list of a deployment

func CmdListDeploymentFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-deployment-file [name] [creator]",
		Short: "list all file names of a deployment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argName := args[0]
			argCreator := args[1]

			params := &types.QueryDeploymentFileNamesRequest{
				SiteName:   argName,
				Creator:    argCreator,
				Pagination: pageReq,
			}

			res, err := queryClient.DeploymentFileNames(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-deployment",
		Short: "list all deployment meta",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDeploymentMetaRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DeploymentMetaAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-deployment name creator",
		Short: "shows a deployment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argName := args[0]
			argCreator := args[1]

			params := &types.QueryGetDeploymentRequest{
				Name:    argName,
				Creator: argCreator,
			}

			res, err := queryClient.Deployment(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
