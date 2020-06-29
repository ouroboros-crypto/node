package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/ouroboros/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group ouroboros queries under a subcommand
	ouroborosQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ouroborosQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdProfile(cdc),
		)...,
	)

	return ouroborosQueryCmd
}

func GetCmdProfile(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "profile [address] [coin]",
		Short: "profile address coin",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]
			coin := args[1]

			_, err := sdk.AccAddressFromBech32(address)

			if err != nil {
				fmt.Printf("Wrong address %s \n", address)
				return nil
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/ouroboros/profile/%s/%s", address, coin), nil)

			if err != nil {
				fmt.Printf("Cannot get profile %s \n", address)
				return nil
			}

			var out types.ProfileResolve

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}
