package cli

import (
	"fmt"

	"ghostcloud/x/ghostcloud/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdUpdateDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update name description domain",
		Short: "Update a deployment.",
		Long:  "This command will overwrite the existing deployment with the new one.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argName := args[0]
			argDescription := args[1]
			argDomain := args[2]
			argWebsitePayload := cmd.Flag(FlagWebsitePayload).Value.String()

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateDeploymentRequest{
				Meta: createMeta(argName, argDescription, argDomain, clientCtx.GetFromAddress().String()),
			}
			if argWebsitePayload != FlagDummyDefault {
				payload, err := createPayload(argWebsitePayload)
				if err != nil {
					return fmt.Errorf("unable to create payload: %v", err)
				}
				msg.Payload = payload
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addUpdateFlags(cmd)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
