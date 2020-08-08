package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/ouroboros/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	ouroborosTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ouroborosTxCmd.AddCommand(flags.PostCommands(
		GetCmdBurnExtra(cdc),
		GetCmdUnburnExtra(cdc),
		GetCmdUpdateRegulation(cdc),
	)...)

	return ouroborosTxCmd
}

func GetCmdBurnExtra(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burn-extra",
		Short: "Burn the extra coins",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgBurnExtraCoins(cliCtx.GetFromAddress())

			err := msg.ValidateBasic()

			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
func GetCmdUnburnExtra(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unburn-extra",
		Short: "Unburn the extra coins",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgUnburnExtraCoins(cliCtx.GetFromAddress())

			err := msg.ValidateBasic()

			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdUpdateRegulation(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update-regulation",
		Short: "Updates the regulation",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			httpClient := http.Client{
				Timeout: 5 * time.Second, // 5 seconds timeout
			}

			resp, err := httpClient.Get("https://api.ouroboros-crypto.com/correction/price")

			if err != nil {
				return err
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)

			// Some problems with parsing the body
			if err != nil {
				return err
			}

			price, isOk := sdk.NewIntFromString(string(body))

			if !isOk {
				return nil
			}

			msg := types.NewMsgUpdateRegulation(cliCtx.GetFromAddress(), price)

			err = msg.ValidateBasic()

			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
