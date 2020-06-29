package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/coins/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	coinsTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	coinsTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateCoin(cdc),
	)...)

	return coinsTxCmd
}

func GetCmdCreateCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "createCoin [name] [symbol] [description] [emission]",
		Short: "Creates a new coin",
		Args:  cobra.ExactArgs(4), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			emission, ok := sdk.NewIntFromString(args[3])

			if !ok {
				panic("wrong emission")
			}

			msg := types.NewMsgCreateCoin(cliCtx.GetFromAddress(), args[0], args[1], args[2], emission, true, []types.CoinBalancePosmining{
				{sdk.NewInt(1), sdk.NewInt(1000000), sdk.NewInt(5000)},
				{sdk.NewInt(1000000), sdk.NewInt(10000000), sdk.NewInt(1000)},
				{sdk.NewInt(100000000), sdk.NewIntWithDecimal(1000000, 6), sdk.NewInt(10000)},
			}, []types.CoinStructurePosmining{
				{sdk.NewInt(1), sdk.NewInt(1000000), sdk.NewInt(131)},
				{sdk.NewInt(1000000), sdk.NewInt(10000000), sdk.NewInt(150)},
				{sdk.NewInt(100000000), sdk.NewIntWithDecimal(1000000, 6), sdk.NewInt(180)},
			})

			err := msg.ValidateBasic()

			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
