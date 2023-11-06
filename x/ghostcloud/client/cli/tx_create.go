package cli

import (
	"fmt"

	"ghostcloud/x/ghostcloud/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func validateCreateParams(argName, argDescription, argDomain, argWebsitePayload string) error {
	err := validateName(argName)
	if err != nil {
		return fmt.Errorf("unable to set name: %v", err)
	}

	err = validateDomain(argDomain)
	if err != nil {
		return fmt.Errorf("unable to set domain: %v", err)
	}

	err = validateDescription(argDescription)
	if err != nil {
		return fmt.Errorf("unable to set description: %v", err)
	}

	err = validateWebsitePayload(argWebsitePayload)
	if err != nil {
		return fmt.Errorf("unable to set website payload: %v", err)
	}
	return nil
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
			argDescription := cmd.Flag(FlagDescription).Value.String()
			argDomain := cmd.Flag(FlagDomain).Value.String()
			err = validateCreateParams(argName, argDescription, argDomain, argWebsitePayload)
			if err != nil {
				return fmt.Errorf("unable to validate params: %v", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			payload, err := createPayload(argWebsitePayload)
			if err != nil {
				return fmt.Errorf("unable to create payload: %v", err)
			}

			meta := types.Meta{
				Creator:     clientCtx.GetFromAddress().String(),
				Name:        argName,
				Description: argDescription,
				Domain:      argDomain,
			}

			msg := &types.MsgCreateDeploymentRequest{
				Meta:    &meta,
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
