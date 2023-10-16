package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"ghostcloud/x/ghostcloud/types"
)

func CmdListDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-deployment",
		Short: "list all deployment",
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

			params := &types.QueryAllDeploymentRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DeploymentAll(cmd.Context(), params)
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
		Use:   "show-deployment [name]",
		Short: "shows a deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argName := args[0]

			params := &types.QueryGetDeploymentRequest{
				Name: argName,
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

func CmdShowDeploymentFileContent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-deployment-file-content [site-name] [creator] [file-name]",
		Short: "shows a deployment file content",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argSiteName := args[0]
			argCreator := args[1]
			argFileName := args[2]

			params := &types.QueryGetDeploymentFileContentRequest{
				SiteName: argSiteName,
				Creator:  argCreator,
				FileName: argFileName,
			}

			res, err := queryClient.DeploymentFileContent(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintBytes(res.Content)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
