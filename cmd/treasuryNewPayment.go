package cmd

import (
	"context"
	"strconv"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// used to construct the action data parameter
type paymentActionParam struct {
	Treasurer    eos.AccountName   `json:"treasurer"`
	RedemptionID uint64            `json:"redemption_id"`
	Amount       eos.Asset         `json:"amount"`
	Notes        map[string]string `json:"notes"`
}

var treasuryNewPaymentCmd = &cobra.Command{
	Use:   "newpayment [redemptionID] [amount] [-n network] [-x trx-id] [-m memo]",
	Short: "treasurer only; creates a payment record against a specific redemption request",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		redemptionID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			panic("Cannot read redemption ID to submit payment against; halting")
		}

		amount, err := eos.NewAssetFromString(args[1])
		if err != nil {
			panic("Unable to read the amount as an Asset object; halting")
		}

		notes := make(map[string]string)
		if len(viper.GetString("create-payment-cmd-network")) > 0 {
			notes["network"] = viper.GetString("create-payment-cmd-network")
		}

		if len(viper.GetString("create-payment-cmd-trx-id")) > 0 {
			notes["trx-id"] = viper.GetString("create-payment-cmd-trx-id")
		}

		if len(viper.GetString("treasury-cmd-memo")) > 0 {
			notes["memo"] = viper.GetString("treasury-cmd-memo")
		}

		action := eos.Action{
			Account: eos.AN(viper.GetString("TreasuryContract")),
			Name:    toActionName("newpayment", "new payment action name"),
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(paymentActionParam{
				Treasurer:    eos.AN(viper.GetString("DAOUser")),
				RedemptionID: redemptionID,
				Amount:       amount,
				Notes:        notes,
			}),
		}

		pushEOSCActions(ctx, getAPI(), &action)

	},
}

func init() {
	treasuryCmd.AddCommand(treasuryNewPaymentCmd)
	treasuryNewPaymentCmd.Flags().StringP("network", "n", "", "network and token used to complete the payment (e.g. BTC, ETH_USDT)")
	treasuryNewPaymentCmd.Flags().StringP("trx-id", "x", "", "transaction ID on the network used to complete the payment")
	// treasuryCreatePaymentCmd.Flags().StringP("memo", "m", "", "memo to be added to the payment record on chain")
}
