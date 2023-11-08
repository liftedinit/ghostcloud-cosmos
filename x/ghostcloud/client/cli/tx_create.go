package cli

import (
	"fmt"

	"ghostcloud/x/ghostcloud/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdCreateDeployment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create name website-payload",
		Short: "Create a new deployment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			argName := args[0]
			argWebsitePayload := args[1]
			argDescription := cmd.Flag(FlagDescription).Value.String()
			argDomain := cmd.Flag(FlagDomain).Value.String()

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			err = validateWebsitePayload(argWebsitePayload)
			if err != nil {
				return fmt.Errorf("unable to validate payload: %v", err)
			}
			payload, err := createPayload(argWebsitePayload)
			if err != nil {
				return fmt.Errorf("unable to create payload: %v", err)
			}

			msg := &types.MsgCreateDeploymentRequest{
				Meta:    createMeta(argName, argDescription, argDomain, clientCtx.GetFromAddress().String()),
				Payload: payload,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addCreateFlags(cmd)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
