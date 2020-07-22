package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/posmining/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	posminingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	posminingTxCmd.AddCommand(flags.PostCommands(
		GetCmdReinvest(cdc),
	)...)

	return posminingTxCmd
}

func GetCmdReinvest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reinvest [coin]",
		Short: "Gets posmining of the desired coin by sending a reinvest transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			coin := coins.Coin{
				Default: args[0] == "ouro",
				Symbol:  args[0],
			}

			msg := types.NewMsgReinvest(cliCtx.GetFromAddress(), coin)

			err := msg.ValidateBasic()

			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
