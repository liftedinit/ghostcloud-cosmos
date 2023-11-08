package cli

import (
	"fmt"
	"strings"

	"ghostcloud/x/ghostcloud/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const FlagFilterBy = "filter-by"
const FlagFilterValue = "filter-value"
const FlagFilterOperator = "filter-operator"

var validFilterByChoices = []string{"creator"}
var validFilterByOperators = []string{"equal", "not_equal"}

func addFilterByFlags(cmd *cobra.Command, filterBy *string) {
	f := cmd.Flags()
	f.StringVarP(filterBy, FlagFilterBy, "", "", "Apply filter to listing (options: creator)")
}

func addFilterValueFlag(cmd *cobra.Command, filterValue *string) {
	f := cmd.Flags()
	f.StringVarP(filterValue, FlagFilterValue, "", "", "The associated value for the filter")
}

func addFilterOperatorFlag(cmd *cobra.Command, filterOperator *string) {
	f := cmd.Flags()
	f.StringVarP(filterOperator, FlagFilterOperator, "", "equal", "The operator for the filter (options: eq, neq, contains, not_contains)")
}

func isValidFilterChoice(choice string, validChoices []string) bool {
	for _, validChoice := range validChoices {
		if choice == validChoice {
			return true
		}
	}
	return false
}

func handleFilterBy(filterBy string, filterByValue string, filterByOperator string) error {
	if filterBy != "" {
		if !isValidFilterChoice(filterBy, validFilterByChoices) {
			return fmt.Errorf("invalid choice for --%s: %s, valid choices are: %v", FlagFilterBy, filterBy, validFilterByChoices)
		}
		if err := handleFilterByValue(filterBy, filterByValue); err != nil {
			return err
		}
		if err := handleFilterByOperator(filterByOperator); err != nil {
			return err
		}

	}
	return nil
}

func handleFilterByValue(filterBy string, filterByValue string) error {
	if filterByValue == "" {
		return fmt.Errorf("invalid value for --%s: %s", FlagFilterValue, filterByValue)
	}

	switch filterBy {
	case "creator":
		if _, err := sdk.AccAddressFromBech32(filterByValue); err != nil {
			return fmt.Errorf("invalid address for --%s: %s: %v", FlagFilterValue, filterByValue, err)
		}
	}
	return nil
}

func handleFilterByOperator(filterByOperator string) error {
	if filterByOperator != "" && !isValidFilterChoice(filterByOperator, validFilterByOperators) {
		return fmt.Errorf("invalid choice for --%s: %s, valid choices are: %v", FlagFilterOperator, filterByOperator, validFilterByOperators)
	}
	return nil
}

func handleFilterFlags(filterBy string, filterByValue string, filterByOperator string) error {
	if err := handleFilterBy(filterBy, filterByValue, filterByOperator); err != nil {
		return err
	}
	return nil
}

func buildFilters(filterBy string, filterByValue string, filterByOperator string) []*types.Filter {
	var filters []*types.Filter
	if filterBy != "" {
		filters = append(filters, &types.Filter{
			Field:    types.Filter_Field(types.Filter_Field_value[strings.ToUpper(filterBy)]),
			Operator: types.Filter_Operator(types.Filter_Operator_value[strings.ToUpper(filterByOperator)]),
			Value:    filterByValue,
		})
	}
	return filters
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group ghostcloud queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListDeployments())
	// this line is used by starport scaffolding # 1

	return cmd
}

func CmdListDeployments() *cobra.Command {
	var filterBy string
	var filterByValue string
	var filterByOperator string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all deployments",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return handleFilterFlags(filterBy, filterByValue, filterByOperator)
		},
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

			params := &types.QueryMetasRequest{
				Pagination: pageReq,
			}

			filters := buildFilters(filterBy, filterByValue, filterByOperator)
			if len(filters) > 0 {
				params.Filters = filters
			}

			res, err := queryClient.Metas(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	addFilterByFlags(cmd, &filterBy)
	addFilterValueFlag(cmd, &filterByValue)
	addFilterOperatorFlag(cmd, &filterByOperator)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
