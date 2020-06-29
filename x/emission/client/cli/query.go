package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/emission/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group emission queries under a subcommand
	emissionQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	emissionQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdGetEmission(queryRoute, cdc),
		)...,
	)

	return emissionQueryCmd
}

func GetCmdGetEmission(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [coin]",
		Short: "get coin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryGetEmission, args[0]), nil)

			if err != nil {
				fmt.Printf("could not get emission\n%s\n", err.Error())

				return nil
			}

			var out types.QueryResGetEmission

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}