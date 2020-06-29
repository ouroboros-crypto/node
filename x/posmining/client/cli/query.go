package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/posmining/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group posmining queries under a subcommand
	posminingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	posminingQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdGetPosmining(queryRoute, cdc),
			GetCmdGetPosminingCoin(queryRoute, cdc),
		)...,
	)

	return posminingQueryCmd
}

// Gets the posmining record
func GetCmdGetPosmining(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [address]",
		Short: "get address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			_, err := sdk.AccAddressFromBech32(address)

			if err != nil {
				fmt.Printf("Wrong address %s \n", address)

				return nil
			}

			var coin = "ouro"
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryGetPosmining, address, coin), nil)

			if err != nil {
				fmt.Printf("Cannot get posmining of %s \n", address)
				fmt.Println(err)

				return nil
			}

			var out types.PosminingResolve

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}

// Gets the posmining record
func GetCmdGetPosminingCoin(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "coin [address] [coin]",
		Short: "coin address coin",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			_, err := sdk.AccAddressFromBech32(address)

			if err != nil {
				fmt.Printf("Wrong address %s \n", address)

				return nil
			}

			var coin string

			if (len(args) == 1) {
				coin = "ouro"
			} else {
				coin = args[1]
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryGetPosmining, address, coin), nil)

			if err != nil {
				fmt.Printf("Cannot get posmining of %s \n", address)
				fmt.Println(err)

				return nil
			}

			var out types.PosminingResolve

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}