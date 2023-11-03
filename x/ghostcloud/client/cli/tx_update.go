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
		Use:   "update name",
		Short: "Update a deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argName := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argDescription := cmd.Flag(FlagDescription).Value.String()
			argDomain := cmd.Flag(FlagDomain).Value.String()
			argWebsitePayload := cmd.Flag(FlagWebsitePayload).Value.String()

			if argDescription == "" && argDomain == "" && argWebsitePayload == "" {
				return fmt.Errorf("at least one of the following flags must be set: --description, --domain, --website-payload")
			}

			meta := types.Meta{
				Creator: clientCtx.GetFromAddress().String(),
				Name:    argName,
			}
			if argDomain != "" {
				meta.Domain = argDomain
			}
			if argDescription != "" {
				meta.Description = argDescription
			}
			msg := &types.MsgUpdateDeploymentRequest{
				Meta: &meta,
			}

			if argWebsitePayload != "" {
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
